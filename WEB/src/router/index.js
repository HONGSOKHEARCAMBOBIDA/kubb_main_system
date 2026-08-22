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
        path: 'semester',
        name: 'ឆមាស',
        component: () => import('@/views/semester/semester.vue'),
        meta: { title: 'ឆមាស', icone: 'Document' },
      },
      {
        path: 'generation',
        name: 'ជំនាន់',
        component: () => import('@/views/generation/generation.vue'),
        meta: { title: 'ជំនាន់', icone: 'Document' },
      },
      {
        path: 'academicshift',
        name: 'វេនសិក្សា',
        component: () => import('@/views/academicsihft/academicshift.vue'),
        meta: { title: 'វេនសិក្សា', icone: 'Document' },
      },
      {
        path: 'academicsection',
        name: 'ក្រុម',
        component: () => import('@/views/academicsection/academicsection.vue'),
        meta: { title: 'ក្រុម', icone: 'Document' },
      },
      {
        path: 'term',
        name: 'វគ្គ',
        component: () => import('@/views/term/term.vue'),
        meta: { title: 'ជំនាន់', icone: 'Document' },
      },
      {
        path: 'faculty',
        name: 'មហាវិទ្យាល័យ',
        component: () => import('@/views/faculty/faculty.vue'),
        meta: { title: 'មហាវិទ្យាល័យ', icone: 'Document' },
      },
      {
        path: 'department',
        name: 'ដេប៉ាតេម៉ង',
        component: () => import('@/views/department/department.vue'),
        meta: { title: 'ដេប៉ាតេម៉ង', icone: 'Document' },
      },
      {
        path: 'major',
        name: 'ជំនាញ',
        component: () => import('@/views/major/major.vue'),
        meta: { title: 'ជំនាញ', icone: 'Document' },
      },
      {
        path: 'academicdegree',
        name: 'ថ្លៃសិក្សា',
        component: () => import('@/views/academicdegree/academicdegree.vue'),
        meta: { title: 'ថ្លៃសិក្សា', icone: 'Document' },
      },
      {
        path: 'subject',
        name: 'មុខវិជ្ជា',
        component: () => import('@/views/subject/subject.vue'),
        meta: { title: 'មុខវិជ្ជា', icone: 'Document' },
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
      {
        path: 'feediscountgroup',
        name: 'ប្រភេទនិស្សិត',
        component: () => import('@/views/feediscountgroup/feediscountgroup.vue'),
        meta: { title: 'ប្រភេទនិស្សិត', icone: 'Document' },
      },
      {
        path: 'schoolarship',
        name: 'អាហារូបករណ៏',
        component: () => import('@/views/schoolarship/schoolarship.vue'),
        meta: { title: 'អាហារូបករណ៏', icone: 'Document' },
      },
      {
        path: 'classcurriculum',
        name: 'បើកថ្នាក់',
        component: () => import('@/views/classcurriculum/classcurriculum.vue'),
        meta: { title: 'បើកថ្នាក់', icone: 'Document' },
      },
      {
        path: 'student',
        name: 'និស្សិត',
        component: () => import('@/views/student/student.vue'),
        meta: { title: 'និស្សិត', icone: 'Document' },
      },
      {
        path: 'admission',
        name: 'និស្សិតបានចុះឈ្មោះសិក្សា',
        component: () => import('@/views/admission/admission.vue'),
        meta: { title: 'និស្សិតបានចុះឈ្មោះសិក្សា', icone: 'Document' },
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
