import { ref, computed } from 'vue'

const user = ref(JSON.parse(localStorage.getItem('loggedInUser') || 'null'))
const isLoggedIn = computed(() => user.value !== null)

function login(userInfo) {
  localStorage.setItem('loggedInUser', JSON.stringify(userInfo))
  user.value = userInfo
}

function updateUser(userData) {
  if (user.value) {
    const updatedUser = { ...user.value, ...userData }
    user.value = updatedUser
    localStorage.setItem('loggedInUser', JSON.stringify(updatedUser))
  }
}

function logout() {
  localStorage.removeItem('loggedInUser')
  user.value = null
}
export function useAuth() {
  return { 
    user, 
    isLoggedIn, 
    login, 
    logout, 
    updateUser,
  }
}