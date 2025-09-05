<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Create New Project</h2>
        <button @click="closeModal" class="close-btn">
          <i class="octicon octicon-x"></i>
        </button>
      </div>
      
      <form @submit.prevent="createProject" class="modal-body">
        <div class="form-group">
          <label for="project-name">Project Name *</label>
          <input
            id="project-name"
            v-model="projectData.name"
            type="text"
            required
            placeholder="my-awesome-project"
            class="form-input"
          />
          <small>Use lowercase letters, numbers, and hyphens only</small>
        </div>
        
        <div class="form-group">
          <label for="project-description">Description</label>
          <textarea
            id="project-description"
            v-model="projectData.description"
            placeholder="Describe what this project does..."
            class="form-textarea"
            rows="3"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label for="project-language">Primary Language *</label>
          <select
            id="project-language"
            v-model="projectData.language"
            required
            class="form-select"
          >
            <option value="">Select a language</option>
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
        </div>
        
        <div class="form-group">
          <label for="project-visibility">Visibility</label>
          <div class="radio-group">
            <label class="radio-option">
              <input
                type="radio"
                v-model="projectData.visibility"
                value="private"
                checked
              />
              <span class="radio-label">
                <i class="octicon octicon-lock"></i>
                Private
              </span>
              <small>Only you can see and access this project</small>
            </label>
            
            <label class="radio-option">
              <input
                type="radio"
                v-model="projectData.visibility"
                value="public"
              />
              <span class="radio-label">
                <i class="octicon octicon-globe"></i>
                Public
              </span>
              <small>Anyone can view, but only you can modify</small>
            </label>
          </div>
        </div>
        
        <div class="form-group">
          <label for="project-tags">Tags</label>
          <input
            id="project-tags"
            v-model="projectData.tags"
            type="text"
            placeholder="api, microservice, ai, ml"
            class="form-input"
          />
          <small>Separate tags with commas</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="projectData.autoGenerateLoader"
              checked
            />
            <span>Auto-generate loader files</span>
          </label>
          <small>Creates language-specific loader files for easy integration</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="projectData.initializeGit"
            />
            <span>Initialize Git repository</span>
          </label>
          <small>Creates a Git repository for version control</small>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="projectData.setupCI"
            />
            <span>Setup CI/CD pipeline</span>
          </label>
          <small>Creates GitHub Actions or GitLab CI configuration</small>
        </div>
      </form>
      
      <div class="modal-footer">
        <button @click="closeModal" class="btn btn-secondary">
          Cancel
        </button>
        <button
          @click="createProject"
          :disabled="!isFormValid || creating"
          class="btn btn-primary"
        >
          <i v-if="creating" class="octicon octicon-sync spin"></i>
          <span v-else>Create Project</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// Props and emits
const emit = defineEmits(['close', 'project-created'])

// Reactive state
const creating = ref(false)
const projectData = ref({
  name: '',
  description: '',
  language: '',
  visibility: 'private',
  tags: '',
  autoGenerateLoader: true,
  initializeGit: false,
  setupCI: false
})

// Computed properties
const isFormValid = computed(() => {
  return projectData.value.name.trim() && 
         projectData.value.language &&
         /^[a-z0-9-]+$/.test(projectData.value.name)
})

// Methods
const closeModal = () => {
  emit('close')
}

const createProject = async () => {
  if (!isFormValid.value) return
  
  creating.value = true
  
  try {
    // Parse tags
    const tags = projectData.value.tags
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0)
    
    const projectPayload = {
      ...projectData.value,
      tags,
      createdAt: new Date().toISOString()
    }
    
    // Here you would send the request to your backend
    const response = await fetch('/api/phantomkit/projects', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('phantomkit_api_key')}`
      },
      body: JSON.stringify(projectPayload)
    })
    
    if (response.ok) {
      const newProject = await response.json()
      emit('project-created', newProject)
    } else {
      const error = await response.json()
      alert(`Error creating project: ${error.message}`)
    }
  } catch (error) {
    console.error('Error creating project:', error)
    alert('Error creating project. Please try again.')
  } finally {
    creating.value = false
  }
}
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
  max-width: 600px;
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

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.radio-option {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  cursor: pointer;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.radio-option:hover {
  border-color: #667eea;
  background: #f8fafc;
}

.radio-option input[type="radio"] {
  margin: 0;
  margin-top: 2px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #374151;
}

.radio-option small {
  display: block;
  margin-top: 4px;
  margin-left: 24px;
  color: #6b7280;
  font-size: 0.75rem;
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

.btn-secondary:hover {
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
  
  .radio-option {
    padding: 12px;
  }
}
</style>
