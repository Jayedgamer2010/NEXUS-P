/**
 * Supabase Auth wrapper for the frontend.
 *
 * Replaces custom PHP auth with Supabase Auth.
 * All auth operations go through Supabase (signUp, signIn, etc.)
 * and the session is managed by Supabase.
 */
import { supabase } from './supabase'
import type { Session, User } from '@supabase/supabase-js'

/**
 * Sign up a new user via Supabase Auth.
 * The database trigger `on_auth_user_created` automatically
 * creates a row in `user_profiles` with the user's data.
 */
export async function signUp(email: string, password: string, username?: string) {
  const { data, error } = await supabase.auth.signUp({
    email,
    password,
    options: {
      data: {
        username: username || email.split('@')[0],
      },
    },
  })
  return { data, error }
}

/**
 * Sign in an existing user via email + password.
 */
export async function signIn(email: string, password: string) {
  const { data, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  })
  return { data, error }
}

/**
 * Sign out the current user.
 */
export async function signOut() {
  const { error } = await supabase.auth.signOut()
  return { error }
}

/**
 * Send a password reset email via Supabase.
 */
export async function resetPassword(email: string) {
  const { data, error } = await supabase.auth.resetPasswordForEmail(email, {
    redirectTo: `${window.location.origin}/auth/reset-password`,
  })
  return { data, error }
}

/**
 * Update the current user's password.
 */
export async function updatePassword(newPassword: string) {
  const { data, error } = await supabase.auth.updateUser({
    password: newPassword,
  })
  return { data, error }
}

/**
 * Update the current user's email.
 */
export async function updateEmail(newEmail: string) {
  const { data, error } = await supabase.auth.updateUser({
    email: newEmail,
  })
  return { data, error }
}

/**
 * Get the currently logged-in session.
 */
export async function getSession(): Promise<Session | null> {
  const { data: { session } } = await supabase.auth.getSession()
  return session
}

/**
 * Get the currently logged-in user.
 */
export async function getUser(): Promise<User | null> {
  const { data: { user } } = await supabase.auth.getUser()
  return user
}

/**
 * Subscribe to auth state changes (login/logout).
 * Returns an unsubscribe function.
 */
export function onAuthStateChange(callback: (event: string, session: Session | null) => void) {
  const { data } = supabase.auth.onAuthStateChange((event, session) => {
    callback(event, session)
  })
  return data.subscription
}

/**
 * Fetch the logged-in user's profile from the `user_profiles` table.
 * This is an app-specific query (not Supabase Auth).
 */
export async function getUserProfile() {
  const user = await getUser()
  if (!user) return { data: null, error: new Error('Not authenticated') }

  const { data, error } = await supabase
    .from('user_profiles')
    .select('*')
    .eq('id', user.id)
    .single()

  return { data, error }
}

/**
 * Update the logged-in user's profile fields.
 * Only fields writable by the RLS policy can be updated.
 */
export async function updateUserProfile(updates: Partial<{
  first_name: string
  last_name: string
  username: string
  avatar: string
  background: string
  company_name: string
  vat_number: string
  address1: string
  address2: string
  city: string
  country: string
  state: string
  postcode: string
}>) {
  const user = await getUser()
  if (!user) return { data: null, error: new Error('Not authenticated') }

  const { data, error } = await supabase
    .from('user_profiles')
    .update({ ...updates, updated_at: new Date().toISOString() })
    .eq('id', user.id)
    .select()
    .single()

  return { data, error }
}
