<script setup>
import { ref, reactive, onMounted, computed, nextTick, watch } from 'vue'
import {
  getAcademicDegree,
  createAcademicDegree,
  updateAcademicDegree,
  toggleAcademicDegree,
} from '../../services/academic_degree.service'
import { getMajorByDepartment } from '../../services/major.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import { getFacultyByProgrammes } from '../../services/faculty.service.js'
import { getDepartmentByFaculty } from '../../services/department.service.js'
import { getAcademics } from '../../services/academic.service'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()

const academicdegree = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const columns = [
  { prop: 'name', label: 'ឈ្មោះកម្រិត', minwidth: 160, slot: 'name' },
  { prop: 'major_name', label: 'ជំនាញ', minwidth: 160, slot: 'major_name' },
  { prop: 'programme_name', label: 'កម្រិតសិក្សា', minwidth: 120 },
  { slot: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 120 },
  { prop: 'monthly_fee', label: 'បង់ប្រចាំខែ', minwidth: 110, slot: 'monthlyFee' },
  { prop: 'quarterly_fee', label: 'បង់ប្រចាំត្រីមាស', minwidth: 120, slot: 'quarterlyFee' },
  { prop: 'semesterly_fee', label: 'បង់ប្រចាំឆមាស', minwidth: 120, slot: 'semesterlyFee' },
  { prop: 'yearly_fee', label: 'បង់ប្រចាំឆ្នាំ', minwidth: 110, slot: 'yearlyFee' },
  { prop: 'active', label: 'ស្ថានភាព', minwidth: 100, slot: 'isActive' },
]

/* ---------------- filters ---------------- */
const filters = reactive({
  programme_id: null,
  faculty_id: null,
  department_id: null,
  major_id: null,
  academic_id: null,
})

const programmesOptions = ref([])
const filterFacultyOptions = ref([])
const filterDepartmentOptions = ref([])
const filterMajorOptions = ref([])
const academicOptions = ref([])

/* ---------------- dialog / form state ---------------- */
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const form = reactive({
  academic_id: null,
  major_id: null,
  name: '',
  monthly_fee: null,
  quarterly_fee: null,
  semesterly_fee: null,
  yearly_fee: null,
  description: '',
})

const formProgramID = ref(null)
const formFacultyID = ref(null)
const formDepartmentID = ref(null)

const formFacultyOptions = ref([])
const formDepartmentOptions = ref([])
const formMajorOptions = ref([])

// While true, the cascade watchers below skip their "reset the child
// field" logic. openEdit sets this before it prefills formProgramID /
// formFacultyID / formDepartmentID, otherwise each watcher immediately
// nulls out the very values we just set.
const isPrefillingForm = ref(false)

const rules = {
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  academic_id: [{ required: true, message: 'សូមជ្រើសរើសឆ្នាំសិក្សា', trigger: 'change' }],
  major_id: [{ required: true, message: 'សូមជ្រើសរើសជំនាញ', trigger: 'change' }],
  monthly_fee: [{ required: true, message: 'សូមបញ្ចូលថ្លៃប្រចាំខែ', trigger: 'blur' }],
  quarterly_fee: [{ required: true, message: 'សូមបញ្ចូលថ្លៃប្រចាំត្រីមាស', trigger: 'blur' }],
  semesterly_fee: [{ required: true, message: 'សូមបញ្ចូលថ្លៃប្រចាំឆមាស', trigger: 'blur' }],
  yearly_fee: [{ required: true, message: 'សូមបញ្ចូលថ្លៃប្រចាំឆ្នាំ', trigger: 'blur' }],
}

/* ---------------- option loaders ---------------- */
async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load programmes')
  }
}

async function fetchAcademicOptions() {
  try {
    const res = await getAcademics()
    academicOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academic years')
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

async function loadDepartmentOption(facultyID) {
  if (!facultyID) return []
  try {
    const res = await getDepartmentByFaculty(facultyID)
    return (res.data.data || []).map((d) => ({ label: d.name, value: d.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load departments')
    return []
  }
}

async function loadMajorOption(departmentID) {
  if (!departmentID) return []
  try {
    const res = await getMajorByDepartment(departmentID)
    return (res.data.data || []).map((m) => ({ label: m.name, value: m.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load majors')
    return []
  }
}

/* ---------------- table data ---------------- */
async function fetchDegrees() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filters.programme_id) params.programme_id = filters.programme_id
    if (filters.faculty_id) params.faculty_id = filters.faculty_id
    if (filters.department_id) params.department_id = filters.department_id
    if (filters.major_id) params.major_id = filters.major_id
    if (filters.academic_id) params.academic_id = filters.academic_id

    const res = await getAcademicDegree(params)
    academicdegree.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

/* ---------------- dialog open/close ---------------- */
function resetForm() {
  form.academic_id = null
  form.major_id = null
  form.name = ''
  form.monthly_fee = null
  form.quarterly_fee = null
  form.semesterly_fee = null
  form.yearly_fee = null
  form.description = ''
  formProgramID.value = null
  formFacultyID.value = null
  formDepartmentID.value = null
  formFacultyOptions.value = []
  formDepartmentOptions.value = []
  formMajorOptions.value = []
}

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតកម្រិតថ្មី'
  resetForm()
  dialogVisible.value = true
}

async function openEdit(row) {
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែកម្រិត'
  resetForm()

  isPrefillingForm.value = true
  try {
    form.name = row.name
    form.description = row.description
    form.monthly_fee = row.monthly_fee
    form.quarterly_fee = row.quarterly_fee
    form.semesterly_fee = row.semesterly_fee
    form.yearly_fee = row.yearly_fee
    form.academic_id = row.academic_id ?? null

    // Programme -> Faculty -> Department -> Major
    formProgramID.value = row.programme_id ?? null
    formFacultyOptions.value = await loadFacultyOption(formProgramID.value)
    formFacultyID.value = row.faculty_id ?? null

    formDepartmentOptions.value = await loadDepartmentOption(formFacultyID.value)
    formDepartmentID.value = row.department_id ?? null

    formMajorOptions.value = await loadMajorOption(formDepartmentID.value)
    form.major_id = row.major_id ?? null

    // Fallback: if the row's major isn't in the loaded options (e.g. the
    // list endpoint doesn't return faculty_id/department_id so the chain
    // above came up empty), inject it directly so the select still shows
    // the current value instead of going blank.
    if (form.major_id && !formMajorOptions.value.some((o) => o.value === form.major_id)) {
      formMajorOptions.value = [
        { label: row.major_name, value: row.major_id },
        ...formMajorOptions.value,
      ]
    }

    await nextTick()
  } finally {
    isPrefillingForm.value = false
  }

  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

/* ---------------- submit / toggle ---------------- */
async function handleSubmit() {
  submitting.value = true
  try {
    const payload = {
      academic_id: form.academic_id,
      major_id: form.major_id,
      name: form.name,
      monthly_fee: form.monthly_fee,
      quarterly_fee: form.quarterly_fee,
      semesterly_fee: form.semesterly_fee,
      yearly_fee: form.yearly_fee,
      description: form.description,
    }
    if (isEditing.value) {
      await updateAcademicDegree(editingUuid.value, payload)
    } else {
      await createAcademicDegree(payload)
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchDegrees()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleAcademicDegree(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchDegrees()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

/* ---------------- filter bar cascade ---------------- */
watch(
  () => filters.programme_id,
  async (newVal) => {
    filters.faculty_id = null
    filters.department_id = null
    filters.major_id = null
    filterDepartmentOptions.value = []
    filterMajorOptions.value = []
    filterFacultyOptions.value = await loadFacultyOption(newVal)
    page.value = 1
    fetchDegrees()
  }
)

watch(
  () => filters.faculty_id,
  async (newVal) => {
    filters.department_id = null
    filters.major_id = null
    filterMajorOptions.value = []
    filterDepartmentOptions.value = await loadDepartmentOption(newVal)
    page.value = 1
    fetchDegrees()
  }
)

watch(
  () => filters.department_id,
  async (newVal) => {
    filters.major_id = null
    filterMajorOptions.value = await loadMajorOption(newVal)
    page.value = 1
    fetchDegrees()
  }
)

watch(
  () => filters.major_id,
  () => {
    page.value = 1
    fetchDegrees()
  }
)

watch(
  () => filters.academic_id,
  () => {
    page.value = 1
    fetchDegrees()
  }
)

/* ---------------- form dialog cascade ---------------- */
// Each of these guards on isPrefillingForm: during openEdit we set the
// ids ourselves and load their option lists explicitly, so the watcher
// must not also reset the child field back to null.
watch(formProgramID, async (newVal) => {
  if (isPrefillingForm.value) return
  formFacultyID.value = null
  formDepartmentID.value = null
  form.major_id = null
  formDepartmentOptions.value = []
  formMajorOptions.value = []
  formFacultyOptions.value = await loadFacultyOption(newVal)
})

watch(formFacultyID, async (newVal) => {
  if (isPrefillingForm.value) return
  formDepartmentID.value = null
  form.major_id = null
  formMajorOptions.value = []
  formDepartmentOptions.value = await loadDepartmentOption(newVal)
})

watch(formDepartmentID, async (newVal) => {
  if (isPrefillingForm.value) return
  form.major_id = null
  formMajorOptions.value = await loadMajorOption(newVal)
})

onMounted(() => {
  fetchProgrammeOptions()
  fetchAcademicOptions()
  fetchDegrees()
})
</script>

<template>
  <div class="academic-degree-page">
    <AppFilterBar
      :fields="[
        { slot: 'program', span: 4 },
        { slot: 'faculty', span: 4 },
        { slot: 'department', span: 4 },
        { slot: 'major', span: 4 },
        { slot: 'academic', span: 4 },
        { slot: 'create', span: 4 },
      ]"
      :action-span="4"
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
          placeholder="ដេប៉ាតេម៉ង់"
          :disabled="!filters.faculty_id"
          clearable
        />
      </template>
      <template #major>
        <AppSelect
          v-model="filters.major_id"
          :options="filterMajorOptions"
          placeholder="ជំនាញ"
          :disabled="!filters.department_id"
          clearable
        />
      </template>
      <template #academic>
        <AppSelect
          v-model="filters.academic_id"
          :options="academicOptions"
          placeholder="ឆ្នាំសិក្សា"
          clearable
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="academicdegree"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchDegrees"
    >
      <template #name="{ row }">
        <el-text tag="b">{{ row.name }}</el-text>
      </template>

      <template #major_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.major_code }} - {{ row.major_name }}
        </el-text>
      </template>

      <template #academic_name="{ row }">
        <el-text tag="b">{{ row.academic_code }}</el-text>
      </template>

      <template #monthlyFee="{ row }">
        <el-text tag="b" style="color: crimson">
          {{ Number(row.monthly_fee).toLocaleString() }} $
        </el-text>
      </template>
      <template #quarterlyFee="{ row }">
        <el-text tag="b" style="color: crimson">
          {{ Number(row.quarterly_fee).toLocaleString() }} $
        </el-text>
      </template>
      <template #semesterlyFee="{ row }">
        <el-text tag="b" style="color: crimson">
          {{ Number(row.semesterly_fee).toLocaleString() }} $
        </el-text>
      </template>
      <template #yearlyFee="{ row }">
        <el-text tag="b" style="color: crimson">
          {{ Number(row.yearly_fee).toLocaleString() }} $
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
        <el-tooltip content="បិទ/បើក" placement="top">
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
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះ"
          clearable
          prop="name"
          label="ឈ្មោះ"
        />

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="form.academic_id"
              :options="academicOptions"
              placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
              clearable
              prop="academic_id"
              label="ឆ្នាំសិក្សា"
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="formProgramID"
              :options="programmesOptions"
              placeholder="ជ្រើសរើសកម្រិតសិក្សា"
              clearable
              label="កម្រិតសិក្សា"
            />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="formFacultyID"
              :options="formFacultyOptions"
              placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
              :disabled="!formProgramID"
              clearable
              label="មហាវិទ្យាល័យ"
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="formDepartmentID"
              :options="formDepartmentOptions"
              placeholder="ជ្រើសរើសដេប៉ាតេម៉ង់"
              :disabled="!formFacultyID"
              clearable
              label="ដេប៉ាតេម៉ង់"
            />
          </el-col>
        </el-row>

        <AppSelect
          v-model="form.major_id"
          :options="formMajorOptions"
          placeholder="ជ្រើសរើសជំនាញ"
          :disabled="!formDepartmentID"
          clearable
          prop="major_id"
          label="ជំនាញ"
        />

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model.number="form.monthly_fee"
              placeholder="ថ្លៃប្រចាំខែ"
              :min="0"
              prop="monthly_fee"
              label="ថ្លៃប្រចាំខែ"
              type="number"
            />
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.quarterly_fee"
              placeholder="ថ្លៃប្រចាំត្រីមាស"
              :min="0"
              prop="quarterly_fee"
              label="ថ្លៃប្រចាំត្រីមាស"
              type="number"
            />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model.number="form.semesterly_fee"
              placeholder="ថ្លៃប្រចាំឆមាស"
              :min="0"
              prop="semesterly_fee"
              label="ថ្លៃប្រចាំឆមាស"
              type="number"
            />
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.yearly_fee"
              placeholder="ថ្លៃប្រចាំឆ្នាំ"
              :min="0"
              prop="yearly_fee"
              label="ថ្លៃប្រចាំឆ្នាំ"
              type="number"
            />
          </el-col>
        </el-row>

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
.academic-degree-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
