<template>
  <Teleport to="body">
    <Transition name="confirm-fade">
      <div v-if="modelValue" class="confirm-overlay" @click.self="onCancel">
        <div class="confirm-dialog" role="alertdialog" aria-modal="true">
          <div class="confirm-dialog__icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.8"/>
              <line x1="12" y1="8" x2="12" y2="13" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
              <circle cx="12" cy="16" r="0.8" fill="currentColor"/>
            </svg>
          </div>

          <h3 class="confirm-dialog__title">{{ title }}</h3>
          <p class="confirm-dialog__message">{{ message }}</p>

          <div class="confirm-dialog__actions">
            <button class="confirm-dialog__btn confirm-dialog__btn--ghost" @click="onCancel">
              {{ cancelLabel }}
            </button>
            <button class="confirm-dialog__btn confirm-dialog__btn--danger" @click="onConfirm">
              {{ confirmLabel }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
export default {
  name: 'ConfirmDialog',
  emits: ['update:modelValue', 'confirm', 'cancel'],

  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: 'Подтвердите действие',
    },
    message: {
      type: String,
      default: '',
    },
    confirmLabel: {
      type: String,
      default: 'Удалить',
    },
    cancelLabel: {
      type: String,
      default: 'Отмена',
    },
  },

  methods: {
    onConfirm() {
      this.$emit('update:modelValue', false)
      this.$emit('confirm')
    },

    onCancel() {
      this.$emit('update:modelValue', false)
      this.$emit('cancel')
    },
  },
}
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(8, 10, 16, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  z-index: 1000;
}

.confirm-dialog {
  width: 100%;
  max-width: 380px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 16px;
  padding: 28px 24px 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.confirm-dialog__icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(224, 92, 92, 0.1);
  color: #e05c5c;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.confirm-dialog__title {
  font-size: 17px;
  font-weight: 700;
  color: #e8eaf0;
  margin-bottom: 8px;
  letter-spacing: -0.01em;
}

.confirm-dialog__message {
  font-size: 14px;
  color: #7b82a0;
  line-height: 1.6;
  margin-bottom: 24px;
  word-break: break-word;
}

.confirm-dialog__actions {
  display: flex;
  gap: 10px;
  width: 100%;
}

.confirm-dialog__btn {
  flex: 1;
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  border: none;
  transition: background 0.15s, opacity 0.15s;
}

.confirm-dialog__btn--ghost {
  background: #1e2236;
  color: #c0c4d6;
  border: 1px solid #2e3347;
}

.confirm-dialog__btn--ghost:hover {
  background: #232840;
  color: #e8eaf0;
}

.confirm-dialog__btn--danger {
  background: #e05c5c;
  color: #fff;
}

.confirm-dialog__btn--danger:hover {
  background: #c84a4a;
}

.confirm-fade-enter-active,
.confirm-fade-leave-active {
  transition: opacity 0.15s ease;
}

.confirm-fade-enter-from,
.confirm-fade-leave-to {
  opacity: 0;
}

.confirm-fade-enter-active .confirm-dialog,
.confirm-fade-leave-active .confirm-dialog {
  transition: transform 0.15s ease;
}

.confirm-fade-enter-from .confirm-dialog,
.confirm-fade-leave-to .confirm-dialog {
  transform: scale(0.96);
}
</style>