<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Upload Code to GitVault</h2>
        <button @click="closeModal" class="close-btn">
          <i class="octicon octicon-x"></i>
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Project Selection -->
        <div class="form-group">
          <label for="upload-project">Target Project *</label>
          <select
            id="upload-project"
            v-model="uploadData.projectId"
            required
            class="form-select"
            @change="onProjectChange"
          >
            <option value="">Select a project</option>
            <option
              v-for="project in availableProjects"
              :key="project.id"
              :value="project.id"
            >
              {{ project.name }} ({{ project.language }})
            </option>
          </select>
          <small>Choose the project to upload code to</small>
        </div>
        
        <!-- File Upload Area -->
        <div class="form-group">
          <label>Files to Upload *</label>
          <div
            class="upload-area"
            :class="{ 'drag-over': isDragOver, 'has-files': files.length > 0 }"
            @drop="onDrop"
            @dragover.prevent="onDragOver"
            @dragleave.prevent="onDragLeave"
            @click="triggerFileInput"
          >
            <input
              ref="fileInput"
              type="file"
              multiple
              @change="onFileSelect"
              style="display: none"
            />
            
            <div v-if="files.length === 0" class="upload-placeholder">
              <i class="octicon octicon-cloud-upload"></i>
              <h3>Drag & Drop Files Here</h3>
              <p>or click to browse files</p>
              <small>Supports: .js, .ts, .py, .go, .rs, .java, .cpp, .c, .php, .rb, .swift, .kt, .scala</small>
            </div>
            
            <div v-else class="file-list">
              <div
                v-for="(file, index) in files"
                :key="index"
                class="file-item"
              >
                <div class="file-info">
                  <i :class="getFileIcon(file.name)"></i>
                  <div class="file-details">
                    <span class="file-name">{{ file.name }}</span>
                    <span class="file-size">{{ formatFileSize(file.size) }}</span>
                  </div>
                </div>
                <button
                  @click.stop="removeFile(index)"
                  class="remove-file-btn"
                  title="Remove file"
                >
                  <i class="octicon octicon-x"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Upload Options -->
        <div class="form-group">
          <label for="upload-script-name">Script Name</label>
          <input
            id="upload-script-name"
            v-model="uploadData.scriptName"
            type="text"
            placeholder="main, utils, helper"
            class="form-input"
          />
          <small>Optional: Custom name for the uploaded script(s)</small>
        </div>
        
        <div class="form-group">
          <label for="upload-language">Language Override</label>
          <select
            id="upload-language"
            v-model="uploadData.languageOverride"
            class="form-select"
          >
            <option value="">Auto-detect from files</option>
            <option value="javascript">JavaScript</option>
            <option value="typescript">TypeScript</option>
            <option value="python">Python</option>
            <option value="go">Go</option>
            <option value="rust">Rust</option>
            <option value="java">Java</option>
            <option value="cpp">C++</option>
            <option value="c">C</option>
            <option value="php">PHP</option>
            <option value="ruby">Ruby</option>
            <option value="swift">Swift</option>
            <option value="kotlin">Kotlin</option>
            <option value="scala">Scala</option>
          </select>
          <small>Override automatic language detection</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="uploadData.generateLoader"
              checked
            />
            <span>Generate loader files</span>
          </label>
          <small>Automatically create language-specific loader files</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="uploadData.overwriteExisting"
            />
            <span>Overwrite existing files</span>
          </label>
          <small>Replace files with the same name</small>
        </div>
        
        <!-- Upload Progress -->
        <div v-if="uploading" class="upload-progress">
          <div class="progress-header">
            <span>Uploading files...</span>
            <span class="progress-text">{{ uploadedCount }}/{{ files.length }}</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{ width: `${uploadProgress}%` }"
            ></div>
          </div>
          <div class="progress-details">
            <span>{{ currentFileName }}</span>
            <span>{{ formatFileSize(uploadedBytes) }} / {{ formatFileSize(totalBytes) }}</span>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="closeModal" class="btn btn-secondary" :disabled="uploading">
          Cancel
        </button>
        <button
          @click="startUpload"
          :disabled="!canUpload || uploading"
          class="btn btn-primary"
        >
          <i v-if="uploading" class="octicon octicon-sync spin"></i>
          <span v-else>Upload {{ files.length }} File{{ files.length !== 1 ? 's' : '' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// Props and emits
const emit = defineEmits(['close', 'upload-complete'])

// Reactive state
const fileInput = ref<HTMLInputElement>()
const uploading = ref(false)
const isDragOver = ref(false)
const files = ref<File[]>([])
const availableProjects = ref([])
const uploadedCount = ref(0)
const uploadedBytes = ref(0)
const totalBytes = ref(0)
const currentFileName = ref('')

const uploadData = ref({
  projectId: '',
  scriptName: '',
  languageOverride: '',
  generateLoader: true,
  overwriteExisting: false
})

// Computed properties
const canUpload = computed(() => {
  return files.value.length > 0 && 
         uploadData.value.projectId &&
         !uploading.value
})

const uploadProgress = computed(() => {
  if (files.value.length === 0) return 0
  return (uploadedCount.value / files.value.length) * 100
})

// Methods
const closeModal = () => {
  if (!uploading.value) {
    emit('close')
  }
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const onFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    addFiles(Array.from(target.files))
  }
}

const onDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

const onDragLeave = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
}

const onDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  
  if (event.dataTransfer?.files) {
    addFiles(Array.from(event.dataTransfer.files))
  }
}

const addFiles = (newFiles: File[]) => {
  const validFiles = newFiles.filter(file => {
    const extension = file.name.split('.').pop()?.toLowerCase()
    const validExtensions = ['js', 'ts', 'py', 'go', 'rs', 'java', 'cpp', 'c', 'php', 'rb', 'swift', 'kt', 'scala']
    return validExtensions.includes(extension || '')
  })
  
  files.value.push(...validFiles)
  calculateTotalBytes()
}

const removeFile = (index: number) => {
  files.value.splice(index, 1)
  calculateTotalBytes()
}

const calculateTotalBytes = () => {
  totalBytes.value = files.value.reduce((total, file) => total + file.size, 0)
}

const onProjectChange = () => {
  // Reset script name when project changes
  uploadData.value.scriptName = ''
}

const startUpload = async () => {
  if (!canUpload.value) return
  
  uploading.value = true
  uploadedCount.value = 0
  uploadedBytes.value = 0
  
  try {
    for (let i = 0; i < files.value.length; i++) {
      const file = files.value[i]
      currentFileName.value = file.name
      
      const formData = new FormData()
      formData.append('file', file)
      formData.append('projectId', uploadData.value.projectId)
      formData.append('scriptName', uploadData.value.scriptName || file.name.split('.')[0])
      formData.append('languageOverride', uploadData.value.languageOverride)
      formData.append('generateLoader', uploadData.value.generateLoader.toString())
      formData.append('overwriteExisting', uploadData.value.overwriteExisting.toString())
      
      const response = await fetch('/api/phantomkit/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('phantomkit_api_key')}`
        },
        body: formData
      })
      
      if (response.ok) {
        uploadedCount.value++
        uploadedBytes.value += file.size
      } else {
        const error = await response.json()
        throw new Error(`Failed to upload ${file.name}: ${error.message}`)
      }
    }
    
    emit('upload-complete')
  } catch (error) {
    console.error('Upload error:', error)
    alert(`Upload failed: ${error.message}`)
  } finally {
    uploading.value = false
  }
}

const getFileIcon = (fileName: string) => {
  const extension = fileName.split('.').pop()?.toLowerCase()
  const icons: Record<string, string> = {
    'js': 'octicon octicon-file-code',
    'ts': 'octicon octicon-file-code',
    'py': 'octicon octicon-file-code',
    'go': 'octicon octicon-file-code',
    'rs': 'octicon octicon-file-code',
    'java': 'octicon octicon-file-code',
    'cpp': 'octicon octicon-file-code',
    'c': 'octicon octicon-file-code',
    'php': 'octicon octicon-file-code',
    'rb': 'octicon octicon-file-code',
    'swift': 'octicon octicon-file-code',
    'kt': 'octicon octicon-file-code',
    'scala': 'octicon octicon-file-code'
  }
  return icons[extension || ''] || 'octicon octicon-file'
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Load available projects
const loadProjects = async () => {
  try {
    const response = await fetch('/api/phantomkit/projects', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('phantomkit_api_key')}`
      }
    })
    
    if (response.ok) {
      availableProjects.value = await response.json()
    }
  } catch (error) {
    console.error('Error loading projects:', error)
  }
}

// Lifecycle
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 12px;
  max-width: 700px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 24px 0 24px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 20px;
}

.modal-header h2 {
  margin: 0;
  color: #111827;
  font-size: 1.5rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #374151;
  font-weight: 500;
  font-size: 0.875rem;
}

.form-input,
.form-select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 0.875rem;
  transition: border-color 0.2s ease;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group small {
  display: block;
  margin-top: 6px;
  color: #6b7280;
  font-size: 0.75rem;
}

.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 12px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f9fafb;
}

.upload-area:hover,
.upload-area.drag-over {
  border-color: #667eea;
  background: #f0f4ff;
}

.upload-area.has-files {
  padding: 20px;
  text-align: left;
}

.upload-placeholder i {
  font-size: 3rem;
  color: #9ca3af;
  margin-bottom: 16px;
}

.upload-placeholder h3 {
  margin: 0 0 8px 0;
  color: #374151;
  font-size: 1.25rem;
}

.upload-placeholder p {
  margin: 0 0 16px 0;
  color: #6b7280;
}

.upload-placeholder small {
  color: #9ca3af;
  font-size: 0.875rem;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-info i {
  font-size: 1.5rem;
  color: #667eea;
}

.file-details {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-weight: 500;
  color: #374151;
}

.file-size {
  font-size: 0.875rem;
  color: #6b7280;
}

.remove-file-btn {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.remove-file-btn:hover {
  background: #fef2f2;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  font-weight: 500;
  color: #374151;
}

.checkbox-label input[type="checkbox"] {
  margin: 0;
  width: 18px;
  height: 18px;
}

.upload-progress {
  background: #f9fafb;
  border-radius: 8px;
  padding: 20px;
  border: 1px solid #e5e7eb;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 500;
  color: #374151;
}

.progress-text {
  color: #667eea;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 12px;
}

.progress-fill {
  height: 100%;
  background: #667eea;
  transition: width 0.3s ease;
}

.progress-details {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  color: #6b7280;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
  border-radius: 0 0 12px 12px;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #5a6fd8;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e7eb;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (max-width: 640px) {
  .modal-content {
    margin: 10px;
    max-height: calc(100vh - 20px);
  }
  
  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 16px;
  }
  
  .upload-area {
    padding: 20px;
  }
  
  .file-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
