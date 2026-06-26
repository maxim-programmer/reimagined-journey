<template>
  <div
    class="dropzone"
    :class="{ 'dropzone--active': isDragging, 'dropzone--error': hasError }"
    @dragenter.prevent="onDragEnter"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
    @click="openFilePicker"
  >
    <input
      ref="fileInput"
      type="file"
      multiple
      accept=".pdf,.docx"
      class="dropzone__input"
      @change="onFileInputChange"
    />

    <div class="dropzone__inner">
      <div class="dropzone__icon" :class="{ 'dropzone__icon--active': isDragging }">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
            d="M24 4L24 32M24 4L16 12M24 4L32 12"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
          <path
            d="M8 34V40C8 41.1 8.9 42 10 42H38C39.1 42 40 41.1 40 40V34"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
          />
        </svg>
      </div>

      <p class="dropzone__title">
        {{ isDragging ? 'Отпустите файлы' : 'Перетащите файлы сюда' }}
      </p>
      <p class="dropzone__subtitle">или <span class="dropzone__link">выберите файлы</span></p>
      <p class="dropzone__hint">PDF и DOCX · не более 20 МБ каждый · несколько файлов одновременно</p>

      <p v-if="hasError" class="dropzone__error-msg">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script>
const MAX_FILE_SIZE = 20 * 1024 * 1024
const ALLOWED_TYPES = {
  'application/pdf': true,
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document': true,
}
const ALLOWED_EXTENSIONS = ['.pdf', '.docx']

export default {
  name: 'DropZone',
  emits: ['files-selected'],

  data() {
    return {
      isDragging: false,
      hasError: false,
      errorMessage: '',
      dragCounter: 0,
    }
  },

  methods: {
    openFilePicker() {
      this.$refs.fileInput.click()
    },

    onDragEnter() {
      this.dragCounter++
      this.isDragging = true
    },

    onDragOver() {
      this.isDragging = true
    },

    onDragLeave() {
      this.dragCounter--
      if (this.dragCounter === 0) {
        this.isDragging = false
      }
    },

    onDrop(event) {
      this.isDragging = false
      this.dragCounter = 0
      const files = Array.from(event.dataTransfer.files)
      this.processFiles(files)
    },

    onFileInputChange(event) {
      const files = Array.from(event.target.files)
      this.processFiles(files)
      event.target.value = ''
    },

    processFiles(files) {
      this.hasError = false
      this.errorMessage = ''

      if (files.length === 0) return

      const valid = []
      const errors = []

      for (const file of files) {
        const ext = '.' + file.name.split('.').pop().toLowerCase()

        if (!ALLOWED_EXTENSIONS.includes(ext) && !ALLOWED_TYPES[file.type]) {
          errors.push(`«${file.name}» — недопустимый формат`)
          continue
        }

        if (file.size > MAX_FILE_SIZE) {
          errors.push(`«${file.name}» — размер превышает 20 МБ`)
          continue
        }

        valid.push(file)
      }

      if (errors.length > 0) {
        this.hasError = true
        this.errorMessage = errors.join('; ')
      }

      if (valid.length > 0) {
        this.$emit('files-selected', valid)
      }
    },
  },
}
</script>

<style scoped>
.dropzone {
  position: relative;
  border: 2px dashed #2e3347;
  border-radius: 16px;
  padding: 56px 40px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
  background: #161a27;
  user-select: none;
}

.dropzone:hover {
  border-color: #4f6ef7;
  background: #181d2e;
}

.dropzone--active {
  border-color: #4f6ef7;
  background: #181d2e;
  box-shadow: 0 0 0 4px rgba(79, 110, 247, 0.12);
}

.dropzone--error {
  border-color: #e05c5c;
}

.dropzone__input {
  display: none;
}

.dropzone__inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  pointer-events: none;
}

.dropzone__icon {
  color: #4f6ef7;
  margin-bottom: 8px;
  transition: transform 0.2s ease;
}

.dropzone__icon--active {
  transform: translateY(-6px);
}

.dropzone__title {
  font-size: 18px;
  font-weight: 600;
  color: #e8eaf0;
  letter-spacing: -0.01em;
}

.dropzone__subtitle {
  font-size: 14px;
  color: #7b82a0;
}

.dropzone__link {
  color: #4f6ef7;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.dropzone__hint {
  font-size: 12px;
  color: #4a5070;
  margin-top: 4px;
}

.dropzone__error-msg {
  font-size: 13px;
  color: #e05c5c;
  margin-top: 8px;
  pointer-events: none;
  max-width: 480px;
}
</style>