<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import type { CampaignAppConfig, CampaignUser } from '../types'
import { describeError, login, register, sendVerifyCode } from '../api'
import { useCampaignMessages } from '../messages'

const props = defineProps<{
  open: boolean
  config: CampaignAppConfig
  campaignSlug: string
  /** 一键拿 Key 时默认注册（新用户转化），否则登录。 */
  initialMode?: 'login' | 'register'
  trialNote?: string
}>()

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'success', payload: { user?: CampaignUser; isNew: boolean }): void
}>()

const t = useCampaignMessages()

type Mode = 'login' | 'register'

const mode = ref<Mode>('login')
const email = ref('')
const password = ref('')
const verifyCode = ref('')
const errorMessage = ref('')
const submitting = ref(false)
const sendingCode = ref(false)
const codeCountdown = ref(0)
let countdownTimer: ReturnType<typeof setInterval> | undefined

const registrationOpen = computed(() => props.config.registration_enabled === true)
const emailVerifyRequired = computed(() => props.config.email_verify_enabled === true)
const captchaEnabled = computed(() => props.config.turnstile_enabled === true)

const title = computed(() =>
  mode.value === 'login' ? t('auth.loginTitle') : t('auth.registerTitle'),
)

watch(
  () => props.open,
  async (open) => {
    if (!open) {
      return
    }
    errorMessage.value = ''
    password.value = ''
    verifyCode.value = ''
    mode.value = props.initialMode === 'register' && registrationOpen.value ? 'register' : 'login'
    await nextTick()
  },
)

function switchMode(next: Mode): void {
  if (next === 'register' && !registrationOpen.value) {
    return
  }
  mode.value = next
  errorMessage.value = ''
}

function stopCountdown(): void {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = undefined
  }
}

async function handleSendCode(): Promise<void> {
  if (!email.value.trim() || sendingCode.value || codeCountdown.value > 0) {
    return
  }
  sendingCode.value = true
  errorMessage.value = ''
  try {
    await sendVerifyCode(email.value.trim())
    codeCountdown.value = 60
    stopCountdown()
    countdownTimer = setInterval(() => {
      codeCountdown.value -= 1
      if (codeCountdown.value <= 0) {
        stopCountdown()
      }
    }, 1000)
  } catch (error) {
    errorMessage.value = describeError(error, t('auth.codeSentFailed'))
  } finally {
    sendingCode.value = false
  }
}

async function handleSubmit(): Promise<void> {
  if (submitting.value) {
    return
  }
  errorMessage.value = ''

  const trimmedEmail = email.value.trim()
  if (!trimmedEmail || !password.value) {
    errorMessage.value = t('auth.emailRequired')
    return
  }
  if (mode.value === 'register' && emailVerifyRequired.value && !verifyCode.value.trim()) {
    errorMessage.value = t('auth.verifyRequired')
    return
  }

  submitting.value = true
  try {
    if (mode.value === 'login') {
      const data = await login(trimmedEmail, password.value)
      emit('success', { user: data.user, isNew: false })
    } else {
      const data = await register({
        email: trimmedEmail,
        password: password.value,
        verify_code: verifyCode.value.trim() || undefined,
        campaign: props.campaignSlug,
      })
      emit('success', { user: data.user, isNew: true })
    }
  } catch (error) {
    errorMessage.value = describeError(
      error,
      mode.value === 'login' ? t('auth.loginFailed') : t('auth.registerFailed'),
    )
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div v-if="open" class="c-modal-mask" @click.self="emit('close')">
    <div class="c-modal" role="dialog" aria-modal="true">
      <div class="flex items-start justify-between gap-4">
        <h3 class="text-lg font-semibold">{{ title }}</h3>
        <button
          type="button"
          class="text-sm"
          style="color: var(--c-text-soft)"
          :aria-label="t('auth.close')"
          @click="emit('close')"
        >
          {{ t('auth.close') }}
        </button>
      </div>

      <div v-if="registrationOpen" class="c-tabs mt-4">
        <button
          class="c-tab"
          type="button"
          :aria-selected="mode === 'login'"
          @click="switchMode('login')"
        >
          {{ t('auth.haveAccount') }}
        </button>
        <button
          class="c-tab"
          type="button"
          :aria-selected="mode === 'register'"
          @click="switchMode('register')"
        >
          {{ t('auth.newUser') }}
        </button>
      </div>
      <p v-else class="c-note mt-4">{{ t('auth.registrationClosed') }}</p>

      <form @submit.prevent="handleSubmit">
        <label class="c-field">
          <span class="c-field-label">{{ t('auth.email') }}</span>
          <input
            v-model="email"
            class="c-input"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
            required
          />
        </label>

        <label class="c-field">
          <span class="c-field-label">{{ t('auth.password') }}</span>
          <input
            v-model="password"
            class="c-input"
            type="password"
            :autocomplete="mode === 'login' ? 'current-password' : 'new-password'"
            :placeholder="t('auth.passwordHint')"
            required
          />
        </label>

        <label v-if="mode === 'register' && emailVerifyRequired" class="c-field">
          <span class="c-field-label">{{ t('auth.verifyCode') }}</span>
          <div class="flex gap-2">
            <input
              v-model="verifyCode"
              class="c-input"
              type="text"
              inputmode="numeric"
              maxlength="6"
              autocomplete="one-time-code"
              :placeholder="t('auth.verifyPlaceholder')"
            />
            <button
              class="c-btn c-btn-ghost"
              type="button"
              :disabled="sendingCode || codeCountdown > 0"
              @click="handleSendCode"
            >
              {{
                codeCountdown > 0
                  ? t('auth.sendIn', { seconds: codeCountdown })
                  : t('auth.send')
              }}
            </button>
          </div>
        </label>

        <p
          v-if="mode === 'register' && trialNote"
          class="mt-3 text-sm"
          style="color: var(--c-text-soft)"
        >
          {{ trialNote }}
        </p>

        <p v-if="captchaEnabled" class="c-note mt-3">{{ t('auth.captchaNotice') }}</p>

        <p v-if="errorMessage" class="c-error">{{ errorMessage }}</p>

        <button class="c-btn c-btn-primary mt-5 w-full" type="submit" :disabled="submitting">
          {{
            submitting
              ? t('auth.processing')
              : mode === 'login'
                ? t('auth.login')
                : t('auth.register')
          }}
        </button>
      </form>
    </div>
  </div>
</template>
