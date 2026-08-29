<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import {
  getSubject,
  createSubject,
  updateSubject,
  toggleSubject,
  creategradecomponent
} from '../../services/subject.service.js'
import { getMajorByDepartment } from '../../services/major.service.js'
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
const subject = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// ---- filter bar state (programme -> faculty -> department -> major cascade) ----
const filters = reactive({
  programme_id: null,
  faculty_id: null,
  department_id: null,
  major_id: null,
})
const programmesOptions = ref([])
const filterFacultyOptions = ref([])
const filterDepartmentOptions = ref([])
const filterMajorOptions = ref([])

// ---- create/edit dialog state (its own separate cascade) ----
const formProgramID = ref(null)
const formFacultyID = ref(null)
const formDepartmentID = ref(null)
const formFacultyOptions = ref([])
const formDepartmentOptions = ref([])
const formMajorOptions = ref([])

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះមុខវិជ្ជា', width: 120 },
  { prop: 'credit', label: 'ក្រេឌីត', width: 100 },
  { prop: 'passing_score', label: 'ពិន្ទុជាប់', width: 100 },
  { prop: 'major_name', label: 'ជំនាញ', slot: 'major_name', width: 160 },
  { prop: 'department_name', label: 'ដេប៉ាតេម៉ង', slot: 'department_name', width: 230 },
  { prop: 'faculty_name', label: 'មហាវិទ្យាល័យ', slot: 'faculty_name', width: 250 },
  { prop: 'programme_name', label: 'កម្រិត', slot: 'programme_name', width: 120 },
  { prop: 'description', label: 'ការពិពណ៌នា', width: 200 },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

const gradecomponentcolmn = [
  { prop: 'name', label: 'ឈ្មោះ', minwidth: 120 },
  { prop: 'weight_percentage',slot:'weight_percentage', label: 'ភាគរយ', minwidth: 120 },
  { prop: 'active',slot:'active', label: 'ស្ថានភាព', minwidth: 120 },
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
  major_id: null,
  credit: null,
  passing_score: null,
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  major_id: [{ required: true, message: 'សូមជ្រើសរើសជំនាញ', trigger: 'change' }],
  credit: [{ required: true, message: 'សូមបញ្ចូលក្រេឌីត', trigger: 'blur' }],
  passing_score: [{ required: true, message: 'សូមបញ្ចូលពិន្ទុជាប់', trigger: 'blur' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

// =====================================================================
// Grade Component dialog state
// =====================================================================
const gradeDialogVisible = ref(false)
const gradeSubmitting = ref(false)
const gradeSubjectID = ref(null)      // int, matches GradeComponentRequestCreate.SubjectID
const gradeSubjectLabel = ref('')     // just for the dialog title
const gradeComponents = reactive([])  // [{ name, weight_percentage }]

const gradeTotalWeight = computed(() =>
  gradeComponents.reduce((sum, g) => sum + (Number(g.weight_percentage) || 0), 0)
)

function makeEmptyGradeRow() {
  return { name: '', weight_percentage: null }
}

function defaultGradeRows() {
  return [
    { name: 'Quiz', weight_percentage: 10 },
    { name: 'Attendance', weight_percentage: 10 },
    { name: 'Assignment', weight_percentage: 10 },
    { name: 'Midterm', weight_percentage: 30 },
    { name: 'Final', weight_percentage: 40 },
  ]
}

function openGradeDialog(row) {
  gradeSubjectID.value = row.id
  gradeSubjectLabel.value = `${row.code} - ${row.name}`
  gradeComponents.splice(0, gradeComponents.length, ...defaultGradeRows())
  gradeDialogVisible.value = true
}

function addGradeRow() {
  gradeComponents.push(makeEmptyGradeRow())
}

function removeGradeRow(index) {
  gradeComponents.splice(index, 1)
  if (gradeComponents.length === 0) {
    gradeComponents.push(makeEmptyGradeRow())
  }
}

function closeGradeDialog() {
  gradeDialogVisible.value = false
}

async function submitGradeComponents() {
  const grade = gradeComponents
    .filter((g) => g.name && g.name.trim() !== '')
    .map((g) => ({
      name: g.name.trim(),
      weight_percentage: Number(g.weight_percentage) || 0,
    }))

  if (grade.length === 0) {
    notify.error('សូមបញ្ចូលធាតុពិន្ទុយ៉ាងហោចណាស់មួយ')
    return
  }

  gradeSubmitting.value = true
  try {
    const payload = {
      subject_id: gradeSubjectID.value,
      grade,
    }
    await creategradecomponent(payload)
    notify.success('រក្សាទុកធាតុពិន្ទុបានជោគជ័យ')
    gradeDialogVisible.value = false
    fetchSubject()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save grade components')
  } finally {
    gradeSubmitting.value = false
  }
}
// =====================================================================

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

async function loadMajorOption(departmentID) {
  if (!departmentID) return []
  try {
    const res = await getMajorByDepartment(departmentID)
    return (res.data.data || []).map((m) => ({
      label: m.name,
      value: m.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load majors')
    return []
  }
}

async function fetchSubject() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.major_id) {
      params.major_id = filters.major_id
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
    const res = await getSubject(params)
    subject.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
    console.log(subject.value)
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតមុខវិជ្ជាថ្មី'
  form.code = ''
  form.name = ''
  form.major_id = null
  form.credit = null
  form.passing_score = null
  form.description = ''
  formProgramID.value = null
  formFacultyID.value = null
  formDepartmentID.value = null
  formFacultyOptions.value = []
  formDepartmentOptions.value = []
  formMajorOptions.value = []
  dialogVisible.value = true
}

async function openEdit(row) {
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែមុខវិជ្ជា'
  form.code = row.code
  form.name = row.name
  form.credit = row.credit
  form.passing_score = row.passing_score
  form.description = row.description

  // Rebuild the cascade so the selects show the right labels
  formProgramID.value = row.programme_id ?? null
  formFacultyOptions.value = await loadFacultyOption(formProgramID.value)
  formFacultyID.value = row.faculty_id ?? null
  formDepartmentOptions.value = await loadDepartmentOption(formFacultyID.value)
  formDepartmentID.value = row.department_id ?? null
  formMajorOptions.value = await loadMajorOption(formDepartmentID.value)
  form.major_id = row.major_id ?? null

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
      major_id: form.major_id,
      credit: form.credit,
      passing_score: form.passing_score,
      description: form.description,
    }
    if (isEditing.value) {
      await updateSubject(editingUuid.value, payload)
    } else {
      await createSubject(payload)
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchSubject()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleSubject(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchSubject()
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
    filters.major_id = null
    filterDepartmentOptions.value = []
    filterMajorOptions.value = []
    filterFacultyOptions.value = await loadFacultyOption(newVal)
    page.value = 1
    fetchSubject()
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
    fetchSubject()
  }
)

watch(
  () => filters.department_id,
  async (newVal) => {
    filters.major_id = null
    filterMajorOptions.value = await loadMajorOption(newVal)
    page.value = 1
    fetchSubject()
  }
)

watch(
  () => filters.major_id,
  () => {
    page.value = 1
    fetchSubject()
  }
)

// ---- form dialog cascade ----
watch(formProgramID, async (newVal) => {
  formFacultyID.value = null
  formDepartmentID.value = null
  form.major_id = null
  formDepartmentOptions.value = []
  formMajorOptions.value = []
  formFacultyOptions.value = await loadFacultyOption(newVal)
})

watch(formFacultyID, async (newVal) => {
  formDepartmentID.value = null
  form.major_id = null
  formMajorOptions.value = []
  formDepartmentOptions.value = await loadDepartmentOption(newVal)
})

watch(formDepartmentID, async (newVal) => {
  form.major_id = null
  formMajorOptions.value = await loadMajorOption(newVal)
})

onMounted(() => {
  fetchProgrammeOptions()
  fetchSubject()
})
</script>

<template>
  <div class="subject-page">
    <AppFilterBar
      :fields="[
        { slot: 'program', span: 3 },
        { slot: 'faculty', span: 3 },
        { slot: 'department', span: 3 },
        { slot: 'major', span: 3 },
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
      <template #major>
        <AppSelect
          v-model="filters.major_id"
          :options="filterMajorOptions"
          placeholder="ជំនាញ"
          :disabled="!filters.department_id"
          clearable
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      expandable
      :data="subject"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchSubject"
    >
      <template #major_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.major_code }} - {{ row.major_name }}
        </el-text>
      </template>

      <template #department_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.department_code }} - {{ row.department_name }}
        </el-text>
      </template>

      <template #faculty_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.faculty_code }} - {{ row.faculty_name }}
        </el-text>
      </template>

      <template #programme_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.programme_name }}
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
        <el-tooltip content="Grade Component" placement="top">
          <AppButton icon="Tickets" circle size="small" type="default" plain @click="openGradeDialog(row)" />
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
         <el-divider content-position="left">Grade Component</el-divider>
         <TableCustom
         expandable
        
                 :data="row.grade_component"
        :columns="gradecomponentcolmn"
        :show-pagination="false"
         >

         <template #weight_percentage="{row}">
          <el-text tag="b" style="color: crimson;">{{ row.weight_percentage }} %</el-text>
         </template>

         <template #active="{row}">
          <el-text>{{ row.active === true ? 'សកម្ម' : 'អសកម្ម' }}</el-text>
         </template>
          
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
        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model="form.code"
              placeholder="បញ្ចូលលេខកូដ"
              clearable
              prop="code"
              label="លេខកូដ"
            />
          </el-col>

          <el-col :span="12">
            <AppInput
              v-model="form.name"
              placeholder="បញ្ចូលឈ្មោះ"
              clearable
              prop="name"
              label="ឈ្មោះមុខវិជ្ជា"
            />
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="formProgramID"
              :options="programmesOptions"
              placeholder="ជ្រើសរើសកម្រិតសិក្សា"
              clearable
              label="កម្រិតសិក្សា"
            />
          </el-col>
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
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="formDepartmentID"
              :options="formDepartmentOptions"
              placeholder="ជ្រើសរើសដេប៉ាតេម៉ង"
              :disabled="!formFacultyID"
              clearable
              label="ដេប៉ាតេម៉ង"
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="form.major_id"
              :options="formMajorOptions"
              placeholder="ជ្រើសរើសជំនាញ"
              :disabled="!formDepartmentID"
              clearable
              prop="major_id"
              label="ជំនាញ"
            />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model.number="form.credit"
              placeholder="បញ្ចូលក្រេឌីត"
              type="number"
              clearable
              prop="credit"
              label="ក្រេឌីត"
            />
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.passing_score"
              placeholder="បញ្ចូលពិន្ទុជាប់"
              type="number"
              clearable
              prop="passing_score"
              label="ពិន្ទុជាប់"
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

    <!-- Grade Component dialog -->
    <AppDialog
      v-model:visible="gradeDialogVisible"
      :title="`Grade Component - ${gradeSubjectLabel}`"
      :showDefaultFooter="false"
      width="58%"
    >
      <div class="grade-component-list">
        <el-row
          v-for="(g, idx) in gradeComponents"
          :key="idx"
          :gutter="12"
          class="grade-component-row"
        >
          <el-col :span="12">
            <AppInput
              v-model="g.name"
              placeholder="ឈ្មោះ (ឧ. Midterm)"
              clearable
              label="ឈ្មោះ"
            />
          </el-col>
          <el-col :span="9">
            <AppInput
              v-model.number="g.weight_percentage"
              placeholder="ភាគរយ"
              type="number"
              clearable
              label="ភាគរយ (%)"
            />
          </el-col>
          <el-col :span="3" class="grade-component-row-actions">
            <AppButton
              icon="Delete"
              circle
              size="small"
              type="danger"
              plain
              @click="removeGradeRow(idx)"
            />
          </el-col>
        </el-row>

        <AppButton type="default" icon="Plus" @click="addGradeRow">បន្ថែមធាតុពិន្ទុ</AppButton>

        <div class="grade-component-total" :class="{ 'is-off': gradeTotalWeight !== 100 }">
          សរុប: {{ gradeTotalWeight }}%
        </div>
      </div>

      <div class="grade-component-footer">
        <AppButton type="default" plain @click="closeGradeDialog">ចាកចេញ</AppButton>
        <AppButton type="default" :loading="gradeSubmitting" @click="submitGradeComponents">
          រក្សាទុក
        </AppButton>
      </div>
    </AppDialog>
  </div>
</template>

<style scoped>
.subject-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.grade-component-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.grade-component-row {
  align-items: center;
}

.grade-component-row-actions {
  display: flex;
  justify-content: center;
}

.grade-component-total {
  font-weight: bold;
  text-align: right;
}

.grade-component-total.is-off {
  color: #d93025;
}

.grade-component-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}
</style>