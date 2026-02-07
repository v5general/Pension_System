import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  id: number
  username: string
  name: string
  role: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>('')
  const isLoggedIn = ref<boolean>(false)

  function setUser(userData: User) {
    user.value = userData
    isLoggedIn.value = true
    // Store user info in localStorage
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function getUser() {
    if (!user.value) {
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        user.value = JSON.parse(storedUser)
        isLoggedIn.value = true
      }
    }
    return user.value
  }

  function logout() {
    user.value = null
    isLoggedIn.value = false
    token.value = ''
    localStorage.removeItem('user')
  }

  function isAdmin() {
    return user.value?.role === 'admin' || user.value?.role === 'operator'
  }

  return {
    user,
    token,
    isLoggedIn,
    setUser,
    getUser,
    logout,
    isAdmin
  }
})
