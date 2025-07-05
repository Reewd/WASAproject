import { ref, computed } from 'vue'
import axios from '../services/axios.js'
import { useAuth } from './useAuth.js'

export function usePhotoUpload() {
  const { user } = useAuth()
  
  const selectedPhoto = ref(null)
  const isUploading = ref(false)
  
  // Create preview URL for selected photo
  const photoPreviewUrl = computed(() => {
    return selectedPhoto.value ? URL.createObjectURL(selectedPhoto.value) : null
  })
  
  // Check if a photo is selected
  const hasSelectedPhoto = computed(() => selectedPhoto.value !== null)
  
  const validateFile = (file) => {
    // Validate file type
    if (!file.type.startsWith('image/')) {
      throw new Error('Please select a valid image file')
    }
    
    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      throw new Error('File size must be less than 10MB')
    }
    
    return true
  }
  
  const handlePhotoSelection = (event) => {
    const file = event.target.files[0]
    if (file) {
      try {
        validateFile(file)
        selectedPhoto.value = file
      } catch (error) {
        alert(error.message)
        return
      }
    }
    // Reset input value to allow selecting the same file again
    event.target.value = ''
  }
  
  const removePhoto = () => {
    if (photoPreviewUrl.value) {
      URL.revokeObjectURL(photoPreviewUrl.value)
    }
    selectedPhoto.value = null
  }
  
  const uploadPhoto = async (photoFile = null) => {
    const fileToUpload = photoFile || selectedPhoto.value
    
    if (!fileToUpload) {
      throw new Error('No photo selected to upload')
    }
    
    isUploading.value = true
    
    const formData = new FormData()
    formData.append('imageFile', fileToUpload)
    
    try {
      const response = await axios.post('/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: user.value.userId,
        },
      })
      return response.data
    } catch (error) {
      console.error('Error uploading photo:', error)
      throw error
    } finally {
      isUploading.value = false
    }
  }
  
  const uploadSelectedPhoto = async () => {
    if (!selectedPhoto.value) {
      return null
    }
    
    try {
      const photoData = await uploadPhoto(selectedPhoto.value)
      removePhoto() // Clear selected photo after successful upload
      return photoData
    } catch (error) {
      console.error('Error uploading selected photo:', error)
      throw error
    }
  }
  
  const resetPhotoState = () => {
    removePhoto()
    isUploading.value = false
  }
  
  return {
    selectedPhoto,
    isUploading,
    photoPreviewUrl,
    hasSelectedPhoto,
    handlePhotoSelection,
    removePhoto,
    uploadPhoto,
    uploadSelectedPhoto,
    resetPhotoState,
    validateFile
  }
}
