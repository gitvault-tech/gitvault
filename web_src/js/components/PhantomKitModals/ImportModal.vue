<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Import Project to GitVault</h2>
        <button @click="closeModal" class="close-btn">
          <i class="octicon octicon-x"></i>
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Import Source Selection -->
        <div class="form-group">
          <label>Import Source *</label>
          <div class="source-tabs">
            <button
              v-for="source in importSources"
              :key="source.id"
              @click="selectedSource = source.id"
              :class="['source-tab', { active: selectedSource === source.id }]"
            >
              <i :class="source.icon"></i>
              {{ source.name }}
            </button>
          </div>
        </div>
        
        <!-- Git Repository Import -->
        <div v-if="selectedSource === 'git'" class="import-section">
          <div class="form-group">
            <label for="git-url">Repository URL *</label>
            <input
              id="git-url"
              v-model="gitData.url"
              type="url"
              placeholder="https://github.com/username/repo.git"
              class="form-input"
            />
            <small>HTTPS or SSH URL for the Git repository</small>
          </div>
          
          <div class="form-group">
            <label for="git-branch">Branch</label>
            <input
              id="git-branch"
              v-model="gitData.branch"
              type="text"
              placeholder="main"
              class="form-input"
            />
            <small>Default: main branch</small>
          </div>
          
          <div class="form-group">
            <label for="git-auth">Authentication</label>
            <select
              id="git-auth"
              v-model="gitData.authType"
              class="form-select"
            >
              <option value="none">No authentication (public repo)</option>
              <option value="ssh">SSH Key</option>
              <option value="token">Personal Access Token</option>
              <option value="username">Username/Password</option>
            </select>
          </div>
          
          <!-- SSH Key Authentication -->
          <div v-if="gitData.authType === 'ssh'" class="auth-section">
            <div class="form-group">
              <label for="ssh-key">SSH Private Key</label>
              <textarea
                id="ssh-key"
                v-model="gitData.sshKey"
                placeholder="-----BEGIN OPENSSH PRIVATE KEY-----"
                class="form-textarea"
                rows="6"
              ></textarea>
              <small>Paste your SSH private key</small>
            </div>
          </div>
          
          <!-- Token Authentication -->
          <div v-if="gitData.authType === 'token'" class="auth-section">
            <div class="form-group">
              <label for="git-token">Personal Access Token</label>
              <input
                id="git-token"
                v-model="gitData.token"
                type="password"
                placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
                class="form-input"
              />
              <small>GitHub, GitLab, or Bitbucket personal access token</small>
            </div>
          </div>
          
          <!-- Username/Password Authentication -->
          <div v-if="gitData.authType === 'username'" class="auth-section">
            <div class="form-group">
              <label for="git-username">Username</label>
              <input
                id="git-username"
                v-model="gitData.username"
                type="text"
                placeholder="your-username"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label for="git-password">Password</label>
              <input
                id="git-password"
                v-model="gitData.password"
                type="password"
                placeholder="your-password"
                class="form-input"
              />
            </div>
          </div>
        </div>
        
        <!-- Archive Import -->
        <div v-if="selectedSource === 'archive'" class="import-section">
          <div class="form-group">
            <label for="archive-url">Archive URL *</label>
            <input
              id="archive-url"
              v-model="archiveData.url"
              type="url"
              placeholder="https://example.com/project.zip"
              class="form-input"
            />
            <small>Direct link to ZIP, TAR, or TAR.GZ archive</small>
          </div>
          
          <div class="form-group">
            <label for="archive-auth">Authentication (if required)</label>
            <input
              id="archive-auth"
              v-model="archiveData.authHeader"
              type="text"
              placeholder="Authorization: Bearer token"
              class="form-input"
            />
            <small>Custom authorization header if needed</small>
          </div>
        </div>
        
        <!-- Local Directory Import -->
        <div v-if="selectedSource === 'local'" class="import-section">
          <div class="form-group">
            <label for="local-path">Directory Path *</label>
            <input
              id="local-path"
              v-model="localData.path"
              type="text"
              placeholder="/path/to/your/project"
              class="form-input"
            />
            <small>Absolute path to the project directory</small>
          </div>
          
          <div class="form-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="localData.includeHidden"
              />
              <span>Include hidden files (.git, .env, etc.)</span>
            </label>
          </div>
        </div>
        
        <!-- Project Configuration -->
        <div class="form-group">
          <label for="import-project-name">Project Name *</label>
          <input
            id="import-project-name"
            v-model="projectConfig.name"
            type="text"
            placeholder="imported-project"
            class="form-input"
          />
          <small>Name for the new GitVault project</small>
        </div>
        
        <div class="form-group">
          <label for="import-project-description">Description</label>
          <textarea
            id="import-project-description"
            v-model="projectConfig.description"
            placeholder="Imported from external source..."
            class="form-textarea"
            rows="3"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label for="import-language">Primary Language</label>
          <select
            id="import-language"
            v-model="projectConfig.language"
            class="form-select"
          >
            <option value="">Auto-detect from source</option>
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
              v-model="projectConfig.generateLoader"
              checked
            />
            <span>Generate loader files</span>
          </label>
          <small>Create language-specific loader files for easy integration</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="projectConfig.initializeGit"
            />
            <span>Initialize Git repository</span>
          </label>
          <small>Create a new Git repository in GitVault</small>
        </div>
        
        <!-- Import Progress -->
        <div v-if="importing" class="import-progress">
          <div class="progress-header">
            <span>Importing project...</span>
            <span class="progress-text">{{ progressStep }}</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{ width: `${importProgress}%` }"
            ></div>
          </div>
          <div class="progress-details">
            <span>{{ currentStep }}</span>
            <span>{{ progressPercentage }}%</span>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="closeModal" class="btn btn-secondary" :disabled="importing">
          Cancel
        </button>
        <button
          @click="startImport"
          :disabled="!canImport || importing"
          class="btn btn-primary"
        >
          <i v-if="importing" class="octicon octicon-sync spin"></i>
          <span v-else>Start Import</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// Props and emits
const emit = defineEmits(['close', 'import-complete'])

// Reactive state
const importing = ref(false)
const selectedSource = ref('git')
const importProgress = ref(0)
const progressStep = ref('')
const currentStep = ref('')
const progressPercentage = ref(0)

const importSources = [
  { id: 'git', name: 'Git Repository', icon: 'octicon octicon-repo' },
  { id: 'archive', name: 'Archive File', icon: 'octicon octicon-file-zip' },
  { id: 'local', name: 'Local Directory', icon: 'octicon octicon-device-desktop' }
]

const gitData = ref({
  url: '',
  branch: 'main',
  authType: 'none',
  sshKey: '',
  token: '',
  username: '',
  password: ''
})

const archiveData = ref({
  url: '',
  authHeader: ''
})

const localData = ref({
  path: '',
  includeHidden: false
})

const projectConfig = ref({
  name: '',
  description: '',
  language: '',
  generateLoader: true,
  initializeGit: false
})

// Computed properties
const canImport = computed(() => {
  if (!projectConfig.value.name.trim()) return false
  
  switch (selectedSource.value) {
    case 'git':
      return gitData.value.url.trim() !== ''
    case 'archive':
      return archiveData.value.url.trim() !== ''
    case 'local':
      return localData.value.path.trim() !== ''
    default:
      return false
  }
})

// Methods
const closeModal = () => {
  if (!importing.value) {
    emit('close')
  }
}

const startImport = async () => {
  if (!canImport.value) return
  
  importing.value = true
  importProgress.value = 0
  progressPercentage.value = 0
  
  try {
    const importPayload = {
      source: selectedSource.value,
      sourceData: getSourceData(),
      projectConfig: projectConfig.value
    }
    
    // Simulate import progress
    await simulateImportProgress()
    
    // Here you would send the request to your backend
    const response = await fetch('/api/phantomkit/import', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('phantomkit_api_key')}`
      },
      body: JSON.stringify(importPayload)
    })
    
    if (response.ok) {
      const importedProject = await response.json()
      emit('import-complete', importedProject)
    } else {
      const error = await response.json()
      throw new Error(`Import failed: ${error.message}`)
    }
  } catch (error) {
    console.error('Import error:', error)
    alert(`Import failed: ${error.message}`)
  } finally {
    importing.value = false
  }
}

const getSourceData = () => {
  switch (selectedSource.value) {
    case 'git':
      return gitData.value
    case 'archive':
      return archiveData.value
    case 'local':
      return localData.value
    default:
      return {}
  }
}

const simulateImportProgress = async () => {
  const steps = [
    'Connecting to source...',
    'Downloading files...',
    'Analyzing project structure...',
    'Detecting language...',
    'Generating loader files...',
    'Setting up project...',
    'Finalizing import...'
  ]
  
  for (let i = 0; i < steps.length; i++) {
    currentStep.value = steps[i]
    progressStep.value = `Step ${i + 1} of ${steps.length}`
    progressPercentage.value = Math.round(((i + 1) / steps.length) * 100)
    importProgress.value = ((i + 1) / steps.length) * 100
    
    // Simulate processing time
    await new Promise(resolve => setTimeout(resolve, 500))
  }
}

// Watch for source changes to reset form data
const resetFormData = () => {
  gitData.value = {
    url: '',
    branch: 'main',
    authType: 'none',
    sshKey: '',
    token: '',
    username: '',
    password: ''
  }
  
  archiveData.value = {
    url: '',
    authHeader: ''
  }
  
  localData.value = {
    path: '',
    includeHidden: false
  }
  
  projectConfig.value = {
    name: '',
    description: '',
    language: '',
    generateLoader: true,
    initializeGit: false
  }
}

// Watch selectedSource changes
import { watch } from 'vue'
watch(selectedSource, resetFormData)
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
.form-textarea,
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
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-group small {
  display: block;
  margin-top: 6px;
  color: #6b7280;
  font-size: 0.75rem;
}

.source-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.source-tab {
  flex: 1;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #6b7280;
}

.source-tab:hover {
  border-color: #667eea;
  background: #f8fafc;
}

.source-tab.active {
  border-color: #667eea;
  background: #f0f4ff;
  color: #667eea;
}

.source-tab i {
  font-size: 1.5rem;
}

.import-section {
  background: #f9fafb;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 24px;
  border: 1px solid #e5e7eb;
}

.auth-section {
  background: white;
  border-radius: 6px;
  padding: 16px;
  margin-top: 16px;
  border: 1px solid #e5e7eb;
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

.import-progress {
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
  
  .source-tabs {
    flex-direction: column;
  }
  
  .source-tab {
    flex: none;
  }
}
</style>
