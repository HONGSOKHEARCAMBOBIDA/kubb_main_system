<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { getMajor, createMajor, updateMajor, toggleMajor } from '../../services/major.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import { getFacultyByProgrammes } from '../../services/faculty.service.js'
import { getDepartmentByFaculty } from '../../services/department.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()
const major = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const durationIntervalOptions = [
  { label: 'ឆ្នាំ', value: 'year' },
  { label: 'ខែ', value: 'month' },
  { label: 'សប្តាហ៍', value: 'week' },
  { label: 'ថ្ងៃ', value: 'day' },
]

// ---- filter bar state (programme -> faculty -> department cascade) ----
const filters = reactive({ programme_id: null, faculty_id: null, department_id: null })
const programmesOptions = ref([])
const filterFacultyOptions = ref([])
const filterDepartmentOptions = ref([])

// ---- create/edit dialog state (its own separate cascade) ----
const formProgramID = ref(null)
const formFacultyID = ref(null)
const formFacultyOptions = ref([])
const formDepartmentOptions = ref([])

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះជំនាញ' },
  { prop: 'duration_period', label: 'រយៈពេលសិក្សា', slot: 'duration_period', width: 140 },
  { prop: 'programme_name', label: 'កម្រិត', slot: 'programme_name', width: 120 },
  { prop: 'faculty_name', label: 'មហាវិទ្យាល័យ', slot: 'faculty_name', width: 220 },
  { prop: 'department_name', label: 'ដេប៉ាតេម៉ង', slot: 'department_name', width: 200 },
  { prop: 'description', label: 'ការពិពណ៌នា',width:300 },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const form = reactive({
  code: '',
  name: '',
  department_id: null,
  description: '',
  duration_period: '',
  duration_interval: 'year',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  department_id: [{ required: true, message: 'សូមជ្រើសរើសដេប៉ាតេម៉ង', trigger: 'change' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
  duration_period: [{ required: true, message: 'សូមបញ្ចូលរយៈពេលសិក្សា', trigger: 'blur' }],
  duration_interval: [{ required: true, message: 'សូមជ្រើសរើសឯកតារយៈពេល', trigger: 'change' }],
}

async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
  }
}

async function loadFacultyOption(programmeID) {
  if (!programmeID) return []
  try {
    const res = await getFacultyByProgrammes(programmeID)
    return (res.data.data || []).map((f) => ({
      label: f.name,
      value: f.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load faculties')
    return []
  }
}

async function loadDepartmentOption(facultyID) {
  if (!facultyID) return []
  try {
    const res = await getDepartmentByFaculty(facultyID)
    return (res.data.data || []).map((d) => ({
      label: d.name,
      value: d.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load departments')
    return []
  }
}

async function fetchMajor() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.department_id) {
      params.department_id = filters.department_id
    }
    if (filters.faculty_id) {
      params.faculty_id = filters.faculty_id
    }
    if (filters.programme_id) {
      params.programme_id = filters.programme_id
    }
    const res = await getMajor(params)
    major.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតជំនាញថ្មី'
  form.code = ''
  form.name = ''
  form.department_id = null
  form.description = ''
  form.duration_period = ''
  form.duration_interval = 'year'
  formProgramID.value = null
  formFacultyID.value = null
  formFacultyOptions.value = []
  formDepartmentOptions.value = []
  dialogVisible.value = true
}

async function openEdit(row) {
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែជំនាញ'
  form.code = row.code
  form.name = row.name
  form.description = row.description
  form.duration_period = row.duration_period
  form.duration_interval = row.duration_interval

  // Rebuild the cascade so the selects show the right labels
  formProgramID.value = row.programme_id ?? null
  formFacultyOptions.value = await loadFacultyOption(formProgramID.value)
  formFacultyID.value = row.faculty_id ?? null
  formDepartmentOptions.value = await loadDepartmentOption(formFacultyID.value)
  form.department_id = row.department_id ?? null

  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    const payload = {
      code: form.code,
      name: form.name,
      department_id: form.department_id,
      description: form.description,
      duration_period: form.duration_period,
      duration_interval: form.duration_interval,
    }
    if (isEditing.value) {
      await updateMajor(editingUuid.value, payload)
    } else {
      await createMajor(payload)
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchMajor()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleMajor(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchMajor()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

// ---- filter bar cascade ----
watch(
  () => filters.programme_id,
  async (newVal) => {
    filters.faculty_id = null
    filters.department_id = null
    filterDepartmentOptions.value = []
    filterFacultyOptions.value = await loadFacultyOption(newVal)
    page.value = 1
    fetchMajor()
  }
)

watch(
  () => filters.faculty_id,
  async (newVal) => {
    filters.department_id = null
    filterDepartmentOptions.value = await loadDepartmentOption(newVal)
    page.value = 1
    fetchMajor()
  }
)

watch(
  () => filters.department_id,
  () => {
    page.value = 1
    fetchMajor()
  }
)

// ---- form dialog cascade ----
watch(formProgramID, async (newVal) => {
  formFacultyID.value = null
  form.department_id = null
  formDepartmentOptions.value = []
  formFacultyOptions.value = await loadFacultyOption(newVal)
})

watch(formFacultyID, async (newVal) => {
  form.department_id = null
  formDepartmentOptions.value = await loadDepartmentOption(newVal)
})

onMounted(() => {
  fetchProgrammeOptions()
  fetchMajor()
})
</script>

<template>
  <div class="major-page">
    <AppFilterBar
      :fields="[
        { slot: 'program', span: 4 },
        { slot: 'faculty', span: 4 },
        { slot: 'department', span: 4 },
        { slot: 'create', span: 3 },
      ]"
      :action-span="3"
    >
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
      <template #department>
        <AppSelect
          v-model="filters.department_id"
          :options="filterDepartmentOptions"
          placeholder="ដេប៉ាតេម៉ង"
          :disabled="!filters.faculty_id"
          clearable
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="major"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchMajor"
    >
  <template #duration_period="{ row }">
       <el-text tag="b" style="color: crimson;">
         {{ row.duration_period }}
        {{ { year: 'ឆ្នាំ', month: 'ខែ', week: 'សប្តាហ៍', day: 'ថ្ងៃ' }[row.duration_interval] || row.duration_interval }}
       </el-text>
      </template>

      <template #programme_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.programme_name }}
        </el-text>
      </template>

      <template #faculty_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.faculty_code }} - {{ row.faculty_name }}
        </el-text>
      </template>

      <template #department_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.department_code }} - {{ row.department_name }}
        </el-text>
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
        <AppInput
          v-model="form.code"
          placeholder="បញ្ចូលលេខកូដ"
          clearable
          prop="code"
          label="លេខកូដ"
        />
        <AppInput
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះ"
          clearable
          prop="name"
          label="ឈ្មោះជំនាញ"
        />
        <AppSelect
          v-model="formProgramID"
          :options="programmesOptions"
          placeholder="ជ្រើសរើសកម្រិតសិក្សា"
          clearable
          label="កម្រិតសិក្សា"
        />
        <AppSelect
          v-model="formFacultyID"
          :options="formFacultyOptions"
          placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
          :disabled="!formProgramID"
          clearable
          label="មហាវិទ្យាល័យ"
        />
        <AppSelect
          v-model="form.department_id"
          :options="formDepartmentOptions"
          placeholder="ជ្រើសរើសដេប៉ាតេម៉ង"
          :disabled="!formFacultyID"
          clearable
          prop="department_id"
          label="ដេប៉ាតេម៉ង"
        />
        <AppInput
           v-model.number="form.duration_period"
          placeholder="បញ្ចូលរយៈពេលសិក្សា"
          type="number"
          clearable
          prop="duration_period"
          label="រយៈពេលសិក្សា"
        />
        <AppSelect
          v-model="form.duration_interval"
          :options="durationIntervalOptions"
          placeholder="ជ្រើសរើសឯកតារយៈពេល"
          prop="duration_interval"
          label="ឯកតារយៈពេល"
        />
        <AppInput
          v-model="form.description"
          placeholder="បញ្ចូលការពិពណ៌នា"
          clearable
          prop="description"
          label="ការពិពណ៌នា"
          type="textarea"
        />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.major-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>