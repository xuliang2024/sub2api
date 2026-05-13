#!/usr/bin/env node

const { execFile } = require('node:child_process');
const { mkdtemp, rm } = require('node:fs/promises');
const { tmpdir } = require('node:os');
const path = require('node:path');
const { promisify } = require('node:util');

const execFileAsync = promisify(execFile);

async function run(command, args, options = {}) {
  try {
    return await execFileAsync(command, args, {
      maxBuffer: 1024 * 1024 * 20,
      ...options,
    });
  } catch (error) {
    const stdout = error.stdout ? `\nstdout:\n${error.stdout}` : '';
    const stderr = error.stderr ? `\nstderr:\n${error.stderr}` : '';
    error.message = `${error.message}${stdout}${stderr}`;
    throw error;
  }
}

function notaryCredentials() {
  const keychainProfile = process.env.NOTARYTOOL_PROFILE || process.env.APPLE_KEYCHAIN_PROFILE;
  if (keychainProfile) {
    return ['--keychain-profile', keychainProfile];
  }

  const appleId = process.env.XC_APPLE_ID || process.env.APPLE_ID;
  const password = process.env.XC_APPLE_APP_SPECIFIC_PASSWORD || process.env.APPLE_APP_SPECIFIC_PASSWORD;
  const teamId = process.env.XC_APPLE_TEAM_ID || process.env.APPLE_TEAM_ID;

  if (!appleId || !password || !teamId) {
    return null;
  }

  return ['--apple-id', appleId, '--password', password, '--team-id', teamId];
}

async function notarizeApp(context = {}) {
  const { electronPlatformName, appOutDir, packager, isFromMain } = context;
  if (!isFromMain && electronPlatformName !== 'darwin') {
    console.log('Skipping notarization for non-macOS target.');
    return;
  }

  const appName = packager?.appInfo?.productFilename || 'Sub2API Desktop';
  const appPath = isFromMain
    ? process.argv.find((arg) => arg.startsWith('--appPath='))?.slice('--appPath='.length)
    : path.join(appOutDir, `${appName}.app`);

  if (!appPath) {
    throw new Error('Missing app path for notarization.');
  }

  const credentials = notaryCredentials();
  if (!credentials) {
    console.log('Apple notarization credentials are not configured; signed app was not notarized.');
    return;
  }

  const tempDir = await mkdtemp(path.join(tmpdir(), 'sub2api-notary-'));
  const zipPath = path.join(tempDir, `${appName}.zip`);

  try {
    console.log(`Preparing ${appName} for Apple notarization...`);
    await run('/usr/bin/ditto', ['-c', '-k', '--sequesterRsrc', '--keepParent', appPath, zipPath]);

    console.log('Submitting app to Apple notary service...');
    const { stdout } = await run('xcrun', [
      'notarytool',
      'submit',
      zipPath,
      ...credentials,
      '--wait',
      '--output-format',
      'json',
    ]);

    const result = JSON.parse(stdout);
    console.log(`Notary submission ${result.id || 'unknown'} finished with status: ${result.status}`);
    if (result.status !== 'Accepted') {
      throw new Error(`Apple notarization failed with status: ${result.status}`);
    }

    console.log('Stapling notarization ticket...');
    await run('xcrun', ['stapler', 'staple', appPath]);
    await run('xcrun', ['stapler', 'validate', appPath]);
    console.log('Apple notarization and stapling completed.');
  } finally {
    await rm(tempDir, { recursive: true, force: true });
  }
}

exports.default = notarizeApp;

if (require.main === module) {
  notarizeApp({ isFromMain: true }).catch((error) => {
    console.error(error.message);
    process.exit(1);
  });
}
