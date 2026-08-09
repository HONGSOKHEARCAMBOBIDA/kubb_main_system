/**
 * Sidebar menu configuration.
 * Add/remove entries here — Sidebar.vue renders this list recursively so
 * nested menus (children) work out of the box.
 *
 * Fields:
 *  - title: label shown in the sidebar
 *  - icon: Element Plus icon component name (auto-registered globally)
 *  - path: route path (omit for a parent-only item with children)
 *  - permission: optional permission string checked via auth store
 *  - children: optional nested items
 */
export const menuConfig = [
  {
    title: 'Dashboard',
    icon: 'Odometer',
    path: '/dashboard',
  },
  {
    title: 'Settings',
    icon: 'Setting',
    children: [
      { title: 'General', icon: 'Tools', path: '/settings' },
      // {title: 'អ្នកប្រេីប្រាស់', icon: 'User',path: '/users'},
      { title: 'សិទ្ធ', icon: 'UserFilled', path: '/role' },

      { title: 'កម្រិតសិក្សា', icon: 'Document', path: '/programmes' },
      { title: 'ឆ្នាំសិក្សា', icon: 'Document', path: '/academic' },
      { title: 'ឆមាស', icon: 'Document', path: '/semester' },
      { title: 'ជំនាន់', icon: 'Document', path: '/generation' },
      { title: 'វគ្គ', icon: 'Document', path: '/term' },
      { title: 'មហាវិទ្យាល័យ', icon: 'Document', path: '/faculty' },
      { title: 'ដេប៉ាតេម៉ង', icon: 'Document', path: '/department' },
      { title: 'សាលា', icon: 'Document', path: '/school' },
      { title: 'សាខាសាលា', icon: 'Document', path: '/campuse' },
      { title: 'អគ្គា', icon: 'Document', path: '/building' },
      { title: 'ជាន់', icon: 'Document', path: '/floor' },
      { title: 'ការិយាល័យ', icon: 'Document', path: '/school_office' },
      { title: 'បន្ទប់', icon: 'Document', path: '/school_room' },
    ],
  },
]
