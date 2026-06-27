<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__logo">
        <svg width="36" height="36" viewBox="0 0 28 28" fill="none">
          <rect width="28" height="28" rx="8" fill="#4f6ef7"/>
          <path d="M8 20V10l6-4 6 4v10" stroke="#fff" stroke-width="1.8" stroke-linejoin="round"/>
          <rect x="11" y="14" width="6" height="6" rx="1" stroke="#fff" stroke-width="1.6"/>
        </svg>
        <span class="auth-card__logo-text">KnowledgeBase</span>
      </div>

      <div class="auth-card__tabs">
        <button
          class="auth-card__tab"
          :class="{ 'auth-card__tab--active': mode === 'login' }"
          @click="switchMode('login')"
        >
          Вход
        </button>
        <button
          class="auth-card__tab"
          :class="{ 'auth-card__tab--active': mode === 'register' }"
          @click="switchMode('register')"
        >
          Регистрация
        </button>
      </div>

      <form class="auth-form" @submit.prevent="onSubmit">
        <div class="auth-form__field">
          <label class="auth-form__label" for="login">Логин</label>
          <input
            id="login"
            v-model="loginVal"
            class="auth-form__input"
            type="text"
            placeholder="Введите логин"
            autocomplete="username"
            :disabled="loading"
          />
        </div>

        <div class="auth-form__field">
          <label class="auth-form__label" for="password">Пароль</label>
          <div class="auth-form__input-wrap">
            <input
              id="password"
              v-model="passwordVal"
              class="auth-form__input"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Введите пароль"
              autocomplete="current-password"
              :disabled="loading"
            />
            <button
              type="button"
              class="auth-form__eye"
              @click="showPassword = !showPassword"
              :title="showPassword ? 'Скрыть' : 'Показать'"
            >
              <svg v-if="!showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
                <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.8"/>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
                <line x1="1" y1="1" x2="23" y2="23" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                <path d="M10.5 6.1A11 11 0 0 1 12 6c7 0 11 6 11 6a20.6 20.6 0 0 1-2.67 3.51" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                <path d="M6.61 6.61A16.05 16.05 0 0 0 1 12s4 6 11 6a10.7 10.7 0 0 0 5.39-1.39" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                <line x1="9.9" y1="9.9" x2="14.12" y2="14.12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
          <p v-if="mode === 'register'" class="auth-form__hint">Минимум 6 символов</p>
        </div>

        <p v-if="error" class="auth-form__error">{{ error }}</p>

        <button
          class="auth-form__submit"
          type="submit"
          :disabled="loading"
        >
          <svg v-if="loading" class="spin" width="16" height="16" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" stroke-dasharray="28 56" stroke-linecap="round"/>
          </svg>
          <span>{{ loading ? 'Подождите…' : (mode === 'login' ? 'Войти' : 'Зарегистрироваться') }}</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { login, register } from '../api/documents.js'

export default {
  name: 'LoginPage',
  emits: ['authenticated'],

  data() {
    return {
      mode: 'login',
      loginVal: '',
      passwordVal: '',
      showPassword: false,
      loading: false,
      error: '',
    }
  },

  methods: {
    switchMode(mode) {
      this.mode = mode
      this.error = ''
    },

    async onSubmit() {
      this.error = ''
      const l = this.loginVal.trim()
      const p = this.passwordVal

      if (!l || !p) {
        this.error = 'Введите логин и пароль'
        return
      }

      this.loading = true
      try {
        if (this.mode === 'login') {
          const data = await login(l, p)
          localStorage.setItem('kb_token', data.token)
          this.$emit('authenticated', data.user)
        } else {
          await register(l, p)
          const data = await login(l, p)
          localStorage.setItem('kb_token', data.token)
          this.$emit('authenticated', data.user)
        }
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  background: #0f1117;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 16px;
  padding: 36px 32px 40px;
}

.auth-card__logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
  justify-content: center;
}

.auth-card__logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #e8eaf0;
  letter-spacing: -0.02em;
}

.auth-card__tabs {
  display: flex;
  gap: 0;
  background: #1e2236;
  border-radius: 10px;
  padding: 3px;
  margin-bottom: 28px;
}

.auth-card__tab {
  flex: 1;
  padding: 8px 0;
  background: none;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  color: #7b82a0;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.auth-card__tab--active {
  background: #2e3347;
  color: #e8eaf0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.auth-form__field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.auth-form__label {
  font-size: 13px;
  font-weight: 500;
  color: #7b82a0;
}

.auth-form__input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.auth-form__input {
  width: 100%;
  background: #1e2236;
  border: 1.5px solid #2e3347;
  border-radius: 10px;
  padding: 11px 14px;
  font-size: 14px;
  color: #e8eaf0;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}

.auth-form__input-wrap .auth-form__input {
  padding-right: 42px;
}

.auth-form__input::placeholder {
  color: #4a5070;
}

.auth-form__input:focus {
  border-color: #4f6ef7;
}

.auth-form__input:disabled {
  opacity: 0.5;
}

.auth-form__eye {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  cursor: pointer;
  color: #4a5070;
  padding: 4px;
  display: flex;
  align-items: center;
  transition: color 0.15s;
}

.auth-form__eye:hover {
  color: #e8eaf0;
}

.auth-form__hint {
  font-size: 11px;
  color: #4a5070;
}

.auth-form__error {
  font-size: 13px;
  color: #e05c5c;
  background: rgba(224, 92, 92, 0.08);
  border: 1px solid rgba(224, 92, 92, 0.2);
  border-radius: 8px;
  padding: 10px 14px;
}

.auth-form__submit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  background: #4f6ef7;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.15s, opacity 0.15s;
  margin-top: 4px;
}

.auth-form__submit:not(:disabled):hover {
  background: #3d5ce8;
}

.auth-form__submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spin {
  animation: spin 0.9s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 480px) {
  .auth-card {
    padding: 28px 20px 32px;
  }
}
</style>