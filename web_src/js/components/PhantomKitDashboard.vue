<template>
  <div class="phantomkit-dashboard">
    <!-- Header Section -->
    <div class="dashboard-header">
      <div class="header-content">
        <h1 class="dashboard-title">
          <i class="octicon octicon-zap"></i>
          PhantomKit Dashboard
        </h1>
        <p class="dashboard-subtitle">
          Secure Code Storage & Runtime Management
        </p>
      </div>
      
      <!-- API Key Management -->
      <div class="api-key-section">
        <div class="api-key-input">
          <input
            v-model="apiKey"
            type="password"
            placeholder="Enter your PhantomKit API Key"
            class="api-key-field"
            @keyup.enter="validateApiKey"
          />
          <button @click="validateApiKey" class="btn btn-primary">
            <i class="octicon octicon-key"></i>
            Connect
          </button>
        </div>
        <div v-if="isConnected" class="connection-status connected">
          <i class="octicon octicon-check-circle"></i>
          Connected to GitVault
        </div>
      </div>
    </div>

    <!-- Main Dashboard Content -->
    <div v-if="isConnected" class="dashboard-content">
      <!-- Quick Actions -->
      <div class="quick-actions">
        <div class="action-card" @click="showNewProjectModal = true">
          <i class="octicon octicon-plus"></i>
          <h3>New Project</h3>
          <p>Create a new PhantomKit project</p>
        </div>
        
        <div class="action-card" @click="showUploadModal = true">
          <i class="octicon octicon-upload"></i>
          <h3>Upload Code</h3>
          <p>Upload scripts and projects</p>
        </div>
        
        <div class="action-card" @click="showImportModal = true">
          <i class="octicon octicon-download"></i>
          <h3>Import Project</h3>
          <p>Import existing projects</p>
        </div>
        
        <div class="action-card" @click="refreshProjects">
          <i class="octicon octicon-sync"></i>
          <h3>Refresh</h3>
          <p>Update project list</p>
        </div>
      </div>

      <!-- Projects Grid -->
      <div class="projects-section">
        <h2>Your Projects</h2>
        <div v-if="loading" class="loading-spinner">
          <i class="octicon octicon-sync spin"></i>
          Loading projects...
        </div>
        
        <div v-else-if="projects.length === 0" class="empty-state">
          <i class="octicon octicon-repo"></i>
          <h3>No projects yet</h3>
          <p>Create your first project to get started with PhantomKit</p>
          <button @click="showNewProjectModal = true" class="btn btn-primary">
            Create Project
          </button>
        </div>
        
        <div v-else class="projects-grid">
          <div
            v-for="project in projects"
            :key="project.id"
            class="project-card"
            @click="selectProject(project)"
          >
            <div class="project-header">
              <div class="project-icon">
                <i :class="getLanguageIcon(project.language)"></i>
              </div>
              <div class="project-info">
                <h3 class="project-name">{{ project.name }}</h3>
                <p class="project-description">{{ project.description }}</p>
              </div>
              <div class="project-status">
                <span :class="['status-badge', project.status]">
                  {{ project.status }}
                </span>
              </div>
            </div>
            
            <div class="project-details">
              <div class="detail-item">
                <i class="octicon octicon-code"></i>
                <span>{{ project.scriptCount }} scripts</span>
              </div>
              <div class="detail-item">
                <i class="octicon octicon-clock"></i>
                <span>Updated {{ formatDate(project.updatedAt) }}</span>
              </div>
            </div>
            
            <div class="project-actions">
              <button @click.stop="viewProject(project)" class="btn btn-sm btn-outline">
                <i class="octicon octicon-eye"></i>
                View
              </button>
              <button @click.stop="runProject(project)" class="btn btn-sm btn-primary">
                <i class="octicon octicon-play"></i>
                Run
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="recent-activity">
        <h2>Recent Activity</h2>
        <div class="activity-list">
          <div
            v-for="activity in recentActivity"
            :key="activity.id"
            class="activity-item"
          >
            <div class="activity-icon">
              <i :class="getActivityIcon(activity.type)"></i>
            </div>
            <div class="activity-content">
              <p class="activity-text">{{ activity.description }}</p>
              <span class="activity-time">{{ formatDate(activity.timestamp) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Connection Required State -->
    <div v-else class="connection-required">
      <div class="connection-content">
        <i class="octicon octicon-lock"></i>
        <h2>Secure Access Required</h2>
        <p>Connect with your PhantomKit API key to access your projects and scripts</p>
        <div class="api-key-help">
          <h4>How to get your API key:</h4>
          <ol>
            <li>Go to your GitVault account settings</li>
            <li>Navigate to API Keys section</li>
            <li>Generate a new PhantomKit API key</li>
            <li>Copy and paste it above</li>
          </ol>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <NewProjectModal
      v-if="showNewProjectModal"
      @close="showNewProjectModal = false"
      @project-created="onProjectCreated"
    />
    
    <UploadModal
      v-if="showUploadModal"
      @close="showUploadModal = false"
      @upload-complete="onUploadComplete"
    />
    
    <ImportModal
      v-if="showImportModal"
      @close="showImportModal = false"
      @import-complete="onImportComplete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import NewProjectModal from './PhantomKitModals/NewProjectModal.vue'
import UploadModal from './PhantomKitModals/UploadModal.vue'
import ImportModal from './PhantomKitModals/ImportModal.vue'

// Reactive state
const apiKey = ref('')
const isConnected = ref(false)
const loading = ref(false)
const projects = ref([])
const recentActivity = ref([])
const showNewProjectModal = ref(false)
const showUploadModal = ref(false)
const showImportModal = ref(false)

// Project interface
interface Project {
  id: string
  name: string
  description: string
  language: string
  status: 'active' | 'inactive' | 'archived'
  scriptCount: number
  updatedAt: string
  createdAt: string
}

// Activity interface
interface Activity {
  id: string
  type: 'upload' | 'run' | 'create' | 'update'
  description: string
  timestamp: string
}

// Methods
const validateApiKey = async () => {
  if (!apiKey.value.trim()) {
    alert('Please enter your API key')
    return
  }

  loading.value = true
  try {
    // Here you would validate the API key with your backend
    const response = await fetch('/api/phantomkit/validate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${apiKey.value}`
      }
    })

    if (response.ok) {
      isConnected.value = true
      localStorage.setItem('phantomkit_api_key', apiKey.value)
      await loadProjects()
      await loadRecentActivity()
    } else {
      alert('Invalid API key. Please check and try again.')
    }
  } catch (error) {
    console.error('Error validating API key:', error)
    alert('Error connecting to GitVault. Please try again.')
  } finally {
    loading.value = false
  }
}

const loadProjects = async () => {
  try {
    const response = await fetch('/api/phantomkit/projects', {
      headers: {
        'Authorization': `Bearer ${apiKey.value}`
      }
    })
    
    if (response.ok) {
      projects.value = await response.json()
    }
  } catch (error) {
    console.error('Error loading projects:', error)
  }
}

const loadRecentActivity = async () => {
  try {
    const response = await fetch('/api/phantomkit/activity', {
      headers: {
        'Authorization': `Bearer ${apiKey.value}`
      }
    })
    
    if (response.ok) {
      recentActivity.value = await response.json()
    }
  } catch (error) {
    console.error('Error loading activity:', error)
  }
}

const refreshProjects = () => {
  loadProjects()
  loadRecentActivity()
}

const selectProject = (project: Project) => {
  // Navigate to project detail view
  window.location.href = `/phantomkit/projects/${project.id}`
}

const viewProject = (project: Project) => {
  selectProject(project)
}

const runProject = async (project: Project) => {
  try {
    const response = await fetch(`/api/phantomkit/projects/${project.id}/run`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${apiKey.value}`
      }
    })
    
    if (response.ok) {
      const result = await response.json()
      alert(`Project ${project.name} executed successfully!\nOutput: ${result.output}`)
    }
  } catch (error) {
    console.error('Error running project:', error)
    alert('Error running project. Please try again.')
  }
}

const onProjectCreated = (project: Project) => {
  projects.value.unshift(project)
  showNewProjectModal.value = false
}

const onUploadComplete = () => {
  refreshProjects()
  showUploadModal.value = false
}

const onImportComplete = () => {
  refreshProjects()
  showImportModal.value = false
}

// Utility functions
const getLanguageIcon = (language: string) => {
  const icons: Record<string, string> = {
    'javascript': 'octicon octicon-file-code',
    'typescript': 'octicon octicon-file-code',
    'python': 'octicon octicon-file-code',
    'go': 'octicon octicon-file-code',
    'rust': 'octicon octicon-file-code',
    'java': 'octicon octicon-file-code'
  }
  return icons[language] || 'octicon octicon-file-code'
}

const getActivityIcon = (type: string) => {
  const icons: Record<string, string> = {
    'upload': 'octicon octicon-upload',
    'run': 'octicon octicon-play',
    'create': 'octicon octicon-plus',
    'update': 'octicon octicon-pencil'
  }
  return icons[type] || 'octicon octicon-dot'
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60))
  
  if (diffInHours < 1) return 'Just now'
  if (diffInHours < 24) return `${diffInHours}h ago`
  if (diffInHours < 168) return `${Math.floor(diffInHours / 24)}d ago`
  return date.toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  // Check for stored API key
  const storedKey = localStorage.getItem('phantomkit_api_key')
  if (storedKey) {
    apiKey.value = storedKey
    validateApiKey()
  }
})
</script>

<style scoped>
.phantomkit-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
  padding: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.header-content h1 {
  margin: 0 0 10px 0;
  font-size: 2.5rem;
  font-weight: 700;
}

.dashboard-subtitle {
  margin: 0;
  font-size: 1.1rem;
  opacity: 0.9;
}

.api-key-section {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.api-key-input {
  display: flex;
  gap: 10px;
}

.api-key-field {
  padding: 12px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  min-width: 300px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.connection-status.connected {
  color: #28a745;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.action-card {
  padding: 30px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.action-card i {
  font-size: 3rem;
  color: #667eea;
  margin-bottom: 20px;
}

.action-card h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.action-card p {
  margin: 0;
  color: #666;
}

.projects-section h2,
.recent-activity h2 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.5rem;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.project-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.project-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.project-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.project-icon i {
  font-size: 2rem;
  color: #667eea;
}

.project-info {
  flex: 1;
}

.project-name {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 1.1rem;
}

.project-description {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
  text-transform: capitalize;
}

.status-badge.active {
  background: #d4edda;
  color: #155724;
}

.status-badge.inactive {
  background: #f8d7da;
  color: #721c24;
}

.status-badge.archived {
  background: #e2e3e5;
  color: #383d41;
}

.project-details {
  display: flex;
  gap: 20px;
  margin-bottom: 15px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #666;
  font-size: 0.9rem;
}

.project-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 5px;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5a6fd8;
}

.btn-outline {
  background: transparent;
  color: #667eea;
  border: 1px solid #667eea;
}

.btn-outline:hover {
  background: #667eea;
  color: white;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
}

.loading-spinner {
  text-align: center;
  padding: 40px;
  color: #666;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #666;
}

.empty-state i {
  font-size: 4rem;
  color: #ddd;
  margin-bottom: 20px;
}

.activity-list {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon i {
  font-size: 1.5rem;
  color: #667eea;
}

.activity-content {
  flex: 1;
}

.activity-text {
  margin: 0 0 5px 0;
  color: #333;
}

.activity-time {
  color: #666;
  font-size: 0.9rem;
}

.connection-required {
  text-align: center;
  padding: 80px 20px;
  color: #666;
}

.connection-content i {
  font-size: 4rem;
  color: #ddd;
  margin-bottom: 20px;
}

.api-key-help {
  text-align: left;
  max-width: 500px;
  margin: 30px auto 0;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.api-key-help h4 {
  margin: 0 0 15px 0;
  color: #333;
}

.api-key-help ol {
  margin: 0;
  padding-left: 20px;
}

.api-key-help li {
  margin-bottom: 8px;
  color: #555;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }
  
  .api-key-input {
    flex-direction: column;
  }
  
  .api-key-field {
    min-width: auto;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .projects-grid {
    grid-template-columns: 1fr;
  }
}
</style>
