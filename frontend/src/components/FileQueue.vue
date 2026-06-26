<template>
  <div v-if="files.length > 0" class="file-queue">
    <div class="file-queue__header">
      <span class="file-queue__label">Очередь загрузки</span>
      <span class="file-queue__count">{{ files.length }} {{ pluralFiles(files.length) }}</span>
    </div>

    <ul class="file-queue__list">
      <li
        v-for="item in files"
        :key="item.id"
        class="file-item"
        :class="`file-item--${item.status}`"
      >
        <div class="file-item__icon">
          <svg v-if="item.status === 'pending'" width="20" height="20" viewBox="0 0 24 24" fill="none">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
            <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round"/>
          </svg>
          <svg v-else-if="item.status === 'uploading'" width="20" height="20" viewBox="0 0 24 24" fill="none" class="spin">
            <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" stroke-dasharray="28 56" stroke-linecap="round"/>
          </svg>
          <svg v-else-if="item.status === 'indexing'" width="20" height="20" viewBox="0 0 24 24" fill="none" class="spin-slow">
            <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2" stroke-dasharray="42 56" stroke-linecap="round"/>
          </svg>
          <svg v-else-if="item.status === 'done'" width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.8"/>
            <polyline points="8 12 11 15 16 9" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else-if="item.status === 'error'" width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.8"/>
            <line x1="12" y1="8" x2="12" y2="13" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="0.8" fill="currentColor"/>
          </svg>
        </div>

        <div class="file-item__body">
          <div class="file-item__top">
            <div class="file-item__info">
              <span class="file-item__name">{{ item.file.name }}</span>
              <span class="file-item__meta">{{ formatSize(item.file.size) }}</span>
            </div>

            <div class="file-item__right">
              <span class="file-item__status-label">
                <span v-if="item.status === 'pending'">Ожидание</span>
                <span v-else-if="item.status === 'uploading'">Загрузка...</span>
                <span v-else-if="item.status === 'indexing'">Индексация...</span>
                <span v-else-if="item.status === 'done'">Готово</span>
                <span v-else-if="item.status === 'error'" :title="item.error">Ошибка</span>
              </span>

              <button
                v-if="item.status === 'pending' || item.status === 'error'"
                class="file-item__remove"
                @click="$emit('remove', item.id)"
                title="Убрать"
              >
                <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                  <line x1="2" y1="2" x2="12" y2="12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                  <line x1="12" y1="2" x2="2" y2="12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
                </svg>
              </button>
            </div>
          </div>

          <div
            v-if="item.status === 'uploading' || item.status === 'indexing' || item.status === 'done' || item.status === 'error'"
            class="file-item__progress-wrap"
          >
            <div class="file-item__progress-track">
              <div
                class="file-item__progress-bar"
                :class="`file-item__progress-bar--${item.status}`"
                :style="{ width: progressWidth(item.status) }"
              ></div>
            </div>
            <span class="file-item__progress-pct">{{ progressPct(item.status) }}</span>
          </div>
        </div>
      </li>
    </ul>

    <div class="file-queue__actions">
      <button
        class="btn btn--primary"
        :disabled="!hasUploadable || isUploading"
        @click="$emit('upload-all')"
      >
        {{ isUploading ? 'Загрузка…' : 'Загрузить все' }}
      </button>
      <button
        class="btn btn--ghost"
        :disabled="isUploading"
        @click="$emit('clear')"
      >
        Очистить
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FileQueue',
  emits: ['remove', 'upload-all', 'clear'],

  props: {
    files: {
      type: Array,
      default: () => [],
    },
    isUploading: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    hasUploadable() {
      return this.files.some((f) => f.status === 'pending')
    },
  },

  methods: {
    progressWidth(status) {
      const map = {
        uploading: '45%',
        indexing: '80%',
        done: '100%',
        error: '100%',
      }
      return map[status] || '0%'
    },

    progressPct(status) {
      const map = {
        uploading: '45%',
        indexing: '80%',
        done: '100%',
        error: '—',
      }
      return map[status] || '0%'
    },

    formatSize(bytes) {
      if (bytes < 1024) return bytes + ' Б'
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' КБ'
      return (bytes / (1024 * 1024)).toFixed(1) + ' МБ'
    },

    pluralFiles(n) {
      if (n % 10 === 1 && n % 100 !== 11) return 'файл'
      if ([2, 3, 4].includes(n % 10) && ![12, 13, 14].includes(n % 100)) return 'файла'
      return 'файлов'
    },
  },
}
</script>

<style scoped>
.file-queue {
  margin-top: 24px;
  background: #161a27;
  border: 1px solid #232840;
  border-radius: 14px;
  overflow: hidden;
}

.file-queue__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid #232840;
}

.file-queue__label {
  font-size: 13px;
  font-weight: 600;
  color: #7b82a0;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.file-queue__count {
  font-size: 13px;
  color: #4a5070;
}

.file-queue__list {
  list-style: none;
  padding: 8px 0;
}

.file-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 20px;
  transition: background 0.15s;
}

.file-item:hover {
  background: #1b2035;
}

.file-item__icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  color: #4a5070;
  margin-top: 2px;
}

.file-item--uploading .file-item__icon {
  color: #4f6ef7;
}

.file-item--indexing .file-item__icon {
  color: #9b6ef7;
}

.file-item--done .file-item__icon {
  color: #3ec97a;
}

.file-item--error .file-item__icon {
  color: #e05c5c;
}

.file-item__body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.file-item__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.file-item__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-item__name {
  font-size: 14px;
  color: #e8eaf0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-item__meta {
  font-size: 12px;
  color: #4a5070;
}

.file-item__right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.file-item__status-label {
  font-size: 12px;
  color: #4a5070;
}

.file-item--uploading .file-item__status-label {
  color: #4f6ef7;
}

.file-item--indexing .file-item__status-label {
  color: #9b6ef7;
}

.file-item--done .file-item__status-label {
  color: #3ec97a;
}

.file-item--error .file-item__status-label {
  color: #e05c5c;
}

.file-item__remove {
  background: none;
  border: none;
  cursor: pointer;
  color: #4a5070;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s, background 0.15s;
}

.file-item__remove:hover {
  color: #e05c5c;
  background: rgba(224, 92, 92, 0.1);
}

.file-item__progress-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-item__progress-track {
  flex: 1;
  height: 4px;
  background: #232840;
  border-radius: 2px;
  overflow: hidden;
}

.file-item__progress-bar {
  height: 100%;
  border-radius: 2px;
  transition: width 0.4s ease;
}

.file-item__progress-bar--uploading {
  background: #4f6ef7;
  animation: shimmer 1.4s ease-in-out infinite;
  background-size: 200% 100%;
  background-image: linear-gradient(
    90deg,
    #4f6ef7 0%,
    #7b96ff 50%,
    #4f6ef7 100%
  );
}

.file-item__progress-bar--indexing {
  background-image: linear-gradient(
    90deg,
    #9b6ef7 0%,
    #bf9fff 50%,
    #9b6ef7 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s ease-in-out infinite;
}

.file-item__progress-bar--done {
  background: #3ec97a;
}

.file-item__progress-bar--error {
  background: #e05c5c;
}

.file-item__progress-pct {
  font-size: 11px;
  color: #4a5070;
  width: 30px;
  text-align: right;
  flex-shrink: 0;
}

.file-queue__actions {
  display: flex;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid #232840;
}

.btn {
  padding: 8px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: background 0.15s, opacity 0.15s;
}

.btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn--primary {
  background: #4f6ef7;
  color: #fff;
}

.btn--primary:not(:disabled):hover {
  background: #3d5ce8;
}

.btn--ghost {
  background: transparent;
  color: #7b82a0;
  border: 1px solid #232840;
}

.btn--ghost:not(:disabled):hover {
  background: #1b2035;
  color: #e8eaf0;
}

.spin {
  animation: spin 0.9s linear infinite;
}

.spin-slow {
  animation: spin 1.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>