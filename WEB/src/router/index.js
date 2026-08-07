import { createRouter, createWebHistory } from 'vue-router'
import { setupGuards } from './guards'

const DashboardLayout = () => import('@/layouts/DashboardLayout.vue')
const AuthLayout = () => import('@/layouts/AuthLayout.vue')

const routes = [
  {
    path: '/',
    component: AuthLayout,
    children: [
      {
        path: '',
        redirect: '/login',
      },
      {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/Login.vue'),
        meta: { title: 'Login', guestOnly: true },
      },
    ],
  },
  {
    path: '/',
    component: DashboardLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: 'Dashboard', icon: 'Odometer' },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/users/Users.vue'),
        meta: { title: 'Users', icon: 'User' },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/Profile.vue'),
        meta: { title: 'Profile', icon: 'UserFilled' },
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/Settings.vue'),
        meta: { title: 'Settings', icon: 'Setting' },
      },
      {
        path: 'role',
        name: 'សិទ្ធ',
        component: () => import('@/views/role/permission.vue'),
        meta: { title: 'សិទ្ធ', icone: 'UserFilled' },
      },
      {
        path: 'author',
        name: 'អ្នកនិពន្ធ',
        component: () => import('@/views/author/author.vue'),
        meta: { title: 'អ្នកនិពន្ធ', icone: 'UserFilled' },
      },
      {
        path: 'category',
        name: 'ប្រភេទសៀវភៅ',
        component: () => import('@/views/category/category.vue'),
        meta: { title: 'category', icone: 'Document' },
      },
      {
        path: 'programmes',
        name: 'កម្រិតសិក្សា',
        component: () => import('@/views/programmes/programmes.vue'),
        meta: { title: 'កម្រិតសិក្សា', icone: 'Document' },
      },
      {
        path: 'academic',
        name: 'ឆ្នាំសិក្សា',
        component: () => import('@/views/academic/academic.vue'),
        meta: { title: 'ឆ្នាំសិក្សា', icone: 'Document' },
      },
      {
        path: 'generation',
        name: 'ជំនាន់',
        component: () => import('@/views/generation/generation.vue'),
        meta: { title: 'ជំនាន់', icone: 'Document' },
      },
      {
        path: 'term',
        name: 'វគ្គ',
        component: () => import('@/views/term/term.vue'),
        meta: { title: 'ជំនាន់', icone: 'Document' },
      },
      {
        path: 'school',
        name: 'សាលា',
        component: () => import('@/views/school/school.vue'),
        meta: { title: 'វគ្គ', icone: 'Document' },
      },
      {
        path: 'campuse',
        name: 'សាខាសាលា',
        component: () => import('@/views/campuse/campuse.vue'),
        meta: { title: 'សាខាសាលា', icone: 'Document' },
      },
      {
        path: 'building',
        name: 'អគ្គា',
        component: () => import('@/views/building/building.vue'),
        meta: { title: 'អគ្គា', icone: 'Document' },
      },
      {
        path: 'floor',
        name: 'ជាន់',
        component: () => import('@/views/floor/floor.vue'),
        meta: { title: 'ជាន់', icone: 'Document' },
      },
      {
        path: 'school_office',
        name: 'ការិយាល័យ',
        component: () => import('@/views/school_office/school_office.vue'),
        meta: { title: 'ការិយាល័យ', icone: 'Document' },
      },
      {
        path: 'school_room',
        name: 'បន្ទប់',
        component: () => import('@/views/school_room/school_room.vue'),
        meta: { title: 'បន្ទប់', icone: 'Document' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/NotFound.vue'),
    meta: { title: 'Page Not Found' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

setupGuards(router)

export default router
