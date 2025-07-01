import { computed } from 'vue'
import { useAuth } from './useAuth.js'

export const useUser = () => {
  const { user, updateUser, getCurrentUsername, getCurrentUserId } = useAuth()

  const getUserId = computed(() => getCurrentUserId())
  const getUsername = computed(() => getCurrentUsername())

  const updateUserData = (userData) => {
    updateUser(userData)
  }

  return { getUserId, getUsername, updateUserData }
}