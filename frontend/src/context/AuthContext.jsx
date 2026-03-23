import React, { createContext, useState, useCallback, useEffect } from 'react'
import { authService } from '../services/authService'

export const AuthContext = createContext()

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [initialized, setInitialized] = useState(false)

  useEffect(() => {
    try {
      const stored = localStorage.getItem('user')
      const storedToken = localStorage.getItem('token')
      if (stored && storedToken) {
        setUser(JSON.parse(stored))
        setToken(storedToken)
      }
    } catch (error) {
      console.error('Erreur lors du chargement du localStorage:', error)
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    }
    setInitialized(true)
  }, [])

  const login = useCallback(async (email, password) => {
    setLoading(true)
    setError(null)
    try {
      const response = await authService.login(email, password)
      const { token: newToken, user: newUser } = response.data
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(newUser))
      setToken(newToken)
      setUser(newUser)
      return newUser
    } catch (err) {
      const errorMsg = err.response?.data?.message || 'Erreur de connexion'
      setError(errorMsg)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  const register = useCallback(async (email, password, secretPhrase) => {
    setLoading(true)
    setError(null)
    try {
      const response = await authService.register(email, password, secretPhrase)
      const { token: newToken, user: newUser } = response.data
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(newUser))
      setToken(newToken)
      setUser(newUser)
      return newUser
    } catch (err) {
      const errorMsg = err.response?.data?.message || "Erreur d'inscription"
      setError(errorMsg)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  const logout = useCallback(() => {
    authService.logout()
    setToken(null)
    setUser(null)
    setError(null)
  }, [])

  if (!initialized) {
    return null
  }

  const value = {
    user,
    token,
    loading,
    error,
    login,
    register,
    logout,
    isAuthenticated: !!token
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
