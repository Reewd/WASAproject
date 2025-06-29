import { ref, computed } from 'vue'

const user = ref(JSON.parse(localStorage.getItem('loggedInUser') || 'null'))
const isLoggedIn = computed(() => user.value !== null)

function login(userInfo) {
  localStorage.setItem('loggedInUser', JSON.stringify(userInfo))
  user.value = userInfo
}

export function useAuth() {
  return { user, isLoggedIn, login }
}
