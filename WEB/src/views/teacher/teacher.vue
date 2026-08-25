<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import {
  getTeacher,
  createTeacher,
  updateTeacher,
  toggleTeacher,
} from '../../services/teacher.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import { getFacultyByProgrammes } from '../../services/faculty.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'
import AppTabs from '../../components/tap/AppTabs.vue'

const notify = useNotification()
const teacher = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const activeTap = ref('general')

const genderOptions = [
  { label: 'ប្រុស', value: 'MALE' },
  { label: 'ស្រី', value: 'FEMALE' },
]

// ---- filter bar state (programme -> faculty cascade) ----
const filters = reactive({ name: '', programme_id: null, faculty_id: null })
const programmesOptions = ref([])
const filterFacultyOptions = ref([])

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 100 },
  { prop: 'name',slot:'name', label: 'ឈ្មោះ' },
  { prop: 'email', label: 'អ៊ីមែល', width: 200 },
  { prop: 'phone', label: 'លេខទូរស័ព្ទ', width: 130 },
  { prop: 'gender', label: 'ភេទ', slot: 'gender', width: 90 },
  { prop: 'date_of_birth', label: 'ថ្ងៃខែឆ្នាំកំណើត', width: 130 },
  { prop: 'place_of_birth', label: 'ទីកន្លែងកំណើត', width: 130 },
  { prop: 'nationality', label: 'សញ្ជាតិ', width: 130 },
   { prop: 'address', label: 'អាស័យដ្ឋាន', width: 130 },
]

const detial = [
  { prop: 'faculty_code', label: 'លេខកូដមហាវិទ្យាល័យ', minwidth: 100 },
  { prop: 'faculty_name', label: 'ឈ្មោះមហាវិទ្យាល័យ', minwidth: 100 },
  { prop: 'programme_name', label: 'កម្រិត', minwidth: 100 },
]

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const form = reactive({
  email: '',
  name: '',
  date_of_birth: '',
  place_of_birth: '',
  gender: '',
  nationality: '',
  address: '',
  phone: '',
  // each row: { programme_id, faculty_id } — programme_id is UI-only,
  // used to filter facultyOptions per row. It is NOT sent to the backend.
  teacher_faculty: [],
})

// parallel array to form.teacher_faculty: facultyOptionsByRow[i] = options for row i
const facultyOptionsByRow = reactive([])

const rules = {
  email: [{ required: true, message: 'សូមបញ្ចូលអ៊ីមែល', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  date_of_birth: [{ required: true, message: 'សូមបញ្ចូលថ្ងៃខែឆ្នាំកំណើត', trigger: 'blur' }],
  place_of_birth: [{ required: true, message: 'សូមបញ្ចូលទីកន្លែងកំណើត', trigger: 'blur' }],
  gender: [{ required: true, message: 'សូមជ្រើសរើសភេទ', trigger: 'change' }],
  nationality: [{ required: true, message: 'សូមបញ្ចូលសញ្ជាតិ', trigger: 'blur' }],
  address: [{ required: true, message: 'សូមបញ្ចូលអាសយដ្ឋាន', trigger: 'blur' }],
  phone: [{ required: true, message: 'សូមបញ្ចូលលេខទូរស័ព្ទ', trigger: 'blur' }],
}

async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
  }
}

async function loadFacultyOption(programmeID) {
  if (!programmeID) return []
  try {
    const res = await getFacultyByProgrammes(programmeID)
    return (res.data.data || []).map((f) => ({ label: f.name, value: f.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load faculties')
    return []
  }
}

async function fetchTeacher() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filters.name) params.name = filters.name
    if (filters.faculty_id) params.faculty_id = filters.faculty_id
    if (filters.programme_id) params.programme_id = filters.programme_id
    const res = await getTeacher(params)
    teacher.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

// ---- per-row cascade helpers ----
function newTeacherFacultyRow() {
  form.teacher_faculty.push({ programme_id: null, faculty_id: null })
  facultyOptionsByRow.push([])
}

function removeTeacherFacultyRow(index) {
  form.teacher_faculty.splice(index, 1)
  facultyOptionsByRow.splice(index, 1)
}

async function onRowProgrammeChange(index) {
  const row = form.teacher_faculty[index]
  row.faculty_id = null
  facultyOptionsByRow[index] = await loadFacultyOption(row.programme_id)
}

function resetForm() {
  form.email = ''
  form.name = ''
  form.date_of_birth = ''
  form.place_of_birth = ''
  form.gender = ''
  form.nationality = ''
  form.address = ''
  form.phone = ''
  form.teacher_faculty = []
  facultyOptionsByRow.splice(0, facultyOptionsByRow.length)
}

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតគ្រូបង្រៀនថ្មី'
  resetForm()
  dialogVisible.value = true
}

async function openEdit(row) {
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែគ្រូបង្រៀន'

  form.email = row.email
  form.name = row.name
  form.date_of_birth = row.date_of_birth
  form.place_of_birth = row.place_of_birth
  form.gender = row.gender
  form.nationality = row.nationality
  form.address = row.address
  form.phone = row.phone

  const teacherFaculty = row.teacher_faculty || []
  form.teacher_faculty = teacherFaculty.map((tf) => ({
    programme_id: tf.programme_id ?? null,
    faculty_id: tf.faculty_id ?? null,
  }))
  facultyOptionsByRow.splice(0, facultyOptionsByRow.length)
  for (const tf of form.teacher_faculty) {
    facultyOptionsByRow.push(tf.programme_id ? await loadFacultyOption(tf.programme_id) : [])
  }

  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    // backend only wants { faculty_id } per row (see TeacherFacultyRequestCreate)
    const payload = {
      email: form.email,
      name: form.name,
      date_of_birth: form.date_of_birth,
      place_of_birth: form.place_of_birth,
      gender: form.gender,
      nationality: form.nationality,
      address: form.address,
      phone: form.phone,
      teacher_faculty: form.teacher_faculty
        .filter((tf) => tf.faculty_id)
        .map((tf) => ({ faculty_id: tf.faculty_id })),
    }

    if (isEditing.value) {
      await updateTeacher(editingUuid.value, payload)
    } else {
      await createTeacher(payload)
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchTeacher()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleTeacher(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchTeacher()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

// ---- filter bar cascade ----
let filterFetchTimer = null
function debouncedFetch() {
  clearTimeout(filterFetchTimer)
  filterFetchTimer = setTimeout(() => {
    page.value = 1
    fetchTeacher()
  }, 300)
}

import { watch } from 'vue'

watch(
  () => filters.programme_id,
  async (newVal) => {
    filters.faculty_id = null
    filterFacultyOptions.value = await loadFacultyOption(newVal)
    page.value = 1
    fetchTeacher()
  }
)

watch(
  () => filters.faculty_id,
  () => {
    page.value = 1
    fetchTeacher()
  }
)

watch(() => filters.name, debouncedFetch)

onMounted(() => {
  fetchProgrammeOptions()
  fetchTeacher()
})
</script>

<template>
  <div class="teacher-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 4 },
        { slot: 'program', span: 4 },
        { slot: 'faculty', span: 4 },
        { slot: 'create', span: 3 },
      ]"
      :action-span="3"
    >
      <template #search>
        <AppInput v-model="filters.name" placeholder="ស្វែងរកតាមឈ្មោះ" clearable />
      </template>
      <template #program>
        <AppSelect
          v-model="filters.programme_id"
          :options="programmesOptions"
          placeholder="កម្រិតសិក្សា"
          clearable
        />
      </template>
      <template #faculty>
        <AppSelect
          v-model="filters.faculty_id"
          :options="filterFacultyOptions"
          placeholder="មហាវិទ្យាល័យ"
          :disabled="!filters.programme_id"
          clearable
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      expandable
      :data="teacher"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchTeacher"
    >
    <template #name="{row}">
      <el-text tag="b" style="color: black;">
        {{ row.name }}
      </el-text>
    </template>
      <template #gender="{ row }">
        {{ { MALE: 'ប្រុស', FEMALE: 'ស្រី' }[row.gender] || row.gender }}
      </template>

      <template #faculty_name="{ row }">
        <el-text tag="b" style="color: darkcyan"
          >{{ row.faculty_code }} - {{ row.faculty_name }}</el-text
        >
      </template>

      <template #programme_name="{ row }">
        <el-text tag="b" style="color: darkcyan">{{ row.programme_name }}</el-text>
      </template>

      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="default" plain @click="openEdit(row)" />
        </el-tooltip>
        <el-tooltip content="បិទ/បេីក" placement="top">
          <AppButton
            :icon="row.active ? 'CircleClose' : 'CircleCheck'"
            circle
            size="small"
            type="default"
            plain
            @click="toggleStatus(row)"
          />
        </el-tooltip>
      </template>
      <template #expand="{ row }">
<el-divider content-position="left">កម្រិតដែលអាចបង្រៀន</el-divider>
<TableCustom
        expandable
        :data="row.teacher_faculty"
        :columns="detial"
        :show-pagination="false"
>

</TableCustom>
      </template>
    </TableCustom>

    <AppDialog v-model:visible="dialogVisible" :title="dialogTitle" :showDefaultFooter="false">
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :loading="submitting"
        :show-actions="true"
        @submit="handleSubmit"
        @reset="closeDialog"
        submitText="រក្សាទុក"
        resetText="ចាកចេញ"
      >
        <AppTabs
          v-model="activeTap"
          :tabs="[
            { name: 'general', label: 'ព័ត៏មានទូទៅ' },
            { name: 'level', label: 'អាចបង្រៀន' },
          ]"
          tab-position="top"
          stretch="true"
        >
          <template #general>
            <el-card>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppInput
                    v-model="form.name"
                    placeholder="បញ្ចូលឈ្មោះ"
                    clearable
                    prop="name"
                    label="ឈ្មោះ"
                  />
                </el-col>
                <el-col :span="12">
                  <AppInput
                    v-model="form.email"
                    placeholder="បញ្ចូលអ៊ីមែល"
                    clearable
                    prop="email"
                    label="អ៊ីមែល"
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppInput
                    v-model="form.date_of_birth"
                    type="date"
                    prop="date_of_birth"
                    label="ថ្ងៃខែឆ្នាំកំណើត"
                  />
                </el-col>
                <el-col :span="12">
                  <AppInput
                    v-model="form.place_of_birth"
                    placeholder="បញ្ចូលទីកន្លែងកំណើត"
                    clearable
                    prop="place_of_birth"
                    label="ទីកន្លែងកំណើត"
                  />
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="form.gender"
                    :options="genderOptions"
                    placeholder="ជ្រើសរើសភេទ"
                    prop="gender"
                    label="ភេទ"
                  />
                </el-col>
                <el-col :span="12">
                  <AppInput
                    v-model="form.nationality"
                    placeholder="បញ្ចូលសញ្ជាតិ"
                    clearable
                    prop="nationality"
                    label="សញ្ជាតិ"
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppInput
                    v-model="form.phone"
                    placeholder="បញ្ចូលលេខទូរស័ព្ទ"
                    clearable
                    prop="phone"
                    label="លេខទូរស័ព្ទ"
                  />
                </el-col>
                <el-col :span="12">
                  <AppInput
                    v-model="form.address"
                    placeholder="បញ្ចូលអាសយដ្ឋាន"
                    clearable
                    prop="address"
                    label="អាសយដ្ឋាន"
                    type="textarea"
                  />
                </el-col>
              </el-row>
            </el-card>
          </template>
          <template #level>
            <el-card>
              <template #header>
                <div class="section-header">
                  <span>កម្រិតដែលអាចបង្រៀនបាន</span>
                  <AppButton
                    type="warning"
                    plain
                    icon="Plus"
                    size="small"
                    @click="newTeacherFacultyRow"
                  >
                    បន្ថែម
                  </AppButton>
                </div>
              </template>

              <el-card
                v-for="(edu, index) in form.teacher_faculty"
                :key="index"
                class="sub-card"
                shadow="never"
              >
                <el-row :gutter="20" align="middle">
                  <el-col :span="11">
                    <AppSelect
                      v-model="edu.programme_id"
                      :options="programmesOptions"
                      placeholder="ជ្រើសរើសកម្រិតសិក្សា"
                      clearable
                      label="កម្រិតសិក្សា"
                      @change="onRowProgrammeChange(index)"
                    />
                  </el-col>
                  <el-col :span="11">
                    <AppSelect
                      v-model="edu.faculty_id"
                      :options="facultyOptionsByRow[index] || []"
                      placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
                      :disabled="!edu.programme_id"
                      clearable
                      label="មហាវិទ្យាល័យ"
                    />
                  </el-col>
                  <el-col :span="2">
                    <AppButton
                      icon="Delete"
                      circle
                      size="small"
                      type="danger"
                      plain
                       :disabled="form.teacher_faculty.length === 1"
                      @click="removeTeacherFacultyRow(index)"
                    />
                  </el-col>
                </el-row>
              </el-card>
            </el-card>
          </template>
        </AppTabs>
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.teacher-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
