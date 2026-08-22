<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { CreateClassCurriculum } from '../../services/classcurriculmn.service.js'
import { getAcademics } from '../../services/academic.service.js'
import { getGenerationByAcademic } from '../../services/generation.service.js'
import { getTermByGeneation } from '../../services/term.service'
import { getSemesterByAcademic } from '../../services/semester.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import { getFacultyByProgrammes } from '../../services/faculty.service.js'
import { getDepartmentByFaculty } from '../../services/department.service.js'
import { getMajorByDepartment } from '../../services/major.service.js'
import { getAcademicshiftByAcademic } from '../../services/academic_shift.service.js'
import { useNotification } from '../../composables/useNotification'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppForm from '@/components/forms/AppForm.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'

const notify = useNotification()

const submitting = ref(false)
const formRef = ref(null)
const dialogVisible = ref(false)
const isEditMode = ref(false)

/* ---------------- top level form (name / term / major) ---------------- */

function emptyDetailRow() {
  return {
    semester_id: null,
    study_year_id: null,
    academic_shift_id: null,
    midterm_date: '',
    final_date: '',
    type_class: ''
  }
}

function defaultForm() {
  return {
    name: '',
    major_id: null,
    term_id: null,
    class_curriclumn_details: [emptyDetailRow()]
  }
}

const form = reactive(defaultForm())

function addDetailRow() {
  form.class_curriclumn_details.push(emptyDetailRow())
}

function removeDetailRow(index) {
  if (form.class_curriclumn_details.length <= 1) return
  form.class_curriclumn_details.splice(index, 1)
}

const rules = {
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះកម្មវិធីសិក្សា', trigger: 'blur' }],
  major_id: [{ required: true, message: 'សូមជ្រើសរើសជំនាញ', trigger: 'change' }],
  term_id: [{ required: true, message: 'សូមជ្រើសរើសវគ្គ', trigger: 'change' }],
}

/* cascading Program -> Faculty -> Department -> Major */
const formProgramID = ref(null)
const formFacultyID = ref(null)
const formDepartmentID = ref(null)
const programmesOptions = ref([])
const formFacultyOptions = ref([])
const formDepartmentOptions = ref([])
const formMajorOptions = ref([])

/* cascading Academic -> Generation -> Term, and Academic -> Semester / Academic Shift */
const formAcademicId = ref(null)
const formGenerationId = ref(null)
const academicOptions = ref([])
const generationOptions = ref([])
const termOptions = ref([])
const semesterOptions = ref([])
const academicShiftOptions = ref([])

const studyyearOption = [
  { label: 'ឆ្នាំទី1', value: 1 },
  { label: 'ឆ្នាំទី2', value: 2 },
  { label: 'ឆ្នាំទី3', value: 3 },
  { label: 'ឆ្នាំទី4', value: 4 },
]

const typeClassOptions = [
  { label: 'ថ្នាក់ផ្ទាល់', value: 'onclass' },
  { label: 'ថ្នាក់អនឡាញ', value: 'online' },
]

/* ---------------- option loaders ---------------- */

async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load programmes')
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

async function fetchAcademicOptions() {
  try {
    const res = await getAcademics()
    academicOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
  }
}

async function loadGenerationOptions(academicId) {
  if (!academicId) return []
  try {
    const res = await getGenerationByAcademic(academicId)
    return (res.data.data || []).map((g) => ({ label: g.name, value: g.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load generations')
    return []
  }
}

async function loadTermOptions(generationId) {
  if (!generationId) return []
  try {
    const res = await getTermByGeneation(generationId)
    return (res.data.data || []).map((t) => ({ label: t.name, value: t.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load terms')
    return []
  }
}

async function loadSemesterOptions(academicId) {
  if (!academicId) return []
  try {
    const res = await getSemesterByAcademic(academicId)
    return (res.data.data || []).map((s) => ({ label: s.name, value: s.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load semesters')
    return []
  }
}

async function loadAcademicShiftOptions(academicId) {
  if (!academicId) return []
  try {
    const res = await getAcademicshiftByAcademic(academicId)
    return (res.data.data || []).map((s) => ({ label: s.name, value: s.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academic shifts')
    return []
  }
}

/* Program -> Faculty -> Department -> Major */
watch(formProgramID, async (newVal) => {
  formFacultyID.value = null
  formDepartmentID.value = null
  formDepartmentOptions.value = []
  formMajorOptions.value = []
  form.major_id = null
  formFacultyOptions.value = await loadFacultyOption(newVal)
})

watch(formFacultyID, async (newVal) => {
  formDepartmentID.value = null
  formMajorOptions.value = []
  form.major_id = null
  formDepartmentOptions.value = await loadDepartmentOption(newVal)
})

watch(formDepartmentID, async (newVal) => {
  form.major_id = null
  formMajorOptions.value = await loadMajorOption(newVal)
})

/* Academic -> Generation, Semester, Academic Shift */
watch(formAcademicId, async (newVal) => {
  formGenerationId.value = null
  form.term_id = null
  generationOptions.value = []
  termOptions.value = []
  semesterOptions.value = []
  academicShiftOptions.value = []

  if (!newVal) return

  generationOptions.value = await loadGenerationOptions(newVal)
  semesterOptions.value = await loadSemesterOptions(newVal)
  academicShiftOptions.value = await loadAcademicShiftOptions(newVal)
})

/* Generation -> Term */
watch(formGenerationId, async (newVal) => {
  form.term_id = null
  termOptions.value = []

  if (!newVal) return

  termOptions.value = await loadTermOptions(newVal)
})

/* ---------------- dialog open / close / submit ---------------- */

function openCreate() {
  isEditMode.value = false
  resetForm()
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  resetForm()
}

async function handleSubmit() {
  submitting.value = true
  try {
    await CreateClassCurriculum(form)
    notify.success('បង្កើតកម្មវិធីសិក្សាបានជោគជ័យ')
    dialogVisible.value = false
    resetForm()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to create')
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  Object.assign(form, defaultForm())

  formProgramID.value = null
  formFacultyID.value = null
  formDepartmentID.value = null
  formFacultyOptions.value = []
  formDepartmentOptions.value = []
  formMajorOptions.value = []

  formAcademicId.value = null
  formGenerationId.value = null
  generationOptions.value = []
  termOptions.value = []
  semesterOptions.value = []
  academicShiftOptions.value = []

  formRef.value?.resetFields?.()
}

onMounted(() => {
  fetchProgrammeOptions()
  fetchAcademicOptions()
})
</script>

<template>
  <AppFilterBar
    :fields="[
      { slot: 'create', span: 4 },
    ]"
  >
    <template #create>
      <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្នាក់ថ្មី</AppButton>
    </template>
  </AppFilterBar>

  <AppDialog
    v-if="dialogVisible"
    v-model:visible="dialogVisible"
    :title="isEditMode ? 'កែប្រែព័ត៌មានកម្មវិធីសិក្សា' : 'បង្កើតកម្មវិធីសិក្សាថ្មី'"
    :showDefaultFooter="false"
    width="60%"
    @close="closeDialog"
  >
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
            v-model="form.name"
            placeholder="បញ្ចូលឈ្មោះកម្មវិធីសិក្សា"
            clearable
            prop="name"
            label="ឈ្មោះកម្មវិធីសិក្សា"
          />
        </el-col>
      </el-row>

      <!-- cascading select: Academic -> Generation -> Term -->
      <el-row :gutter="20">
        <el-col :span="8">
          <AppSelect
            v-model="formAcademicId"
            :options="academicOptions"
            placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
            clearable
            label="ឆ្នាំសិក្សា"
          />
        </el-col>
        <el-col :span="8">
          <AppSelect
            v-model="formGenerationId"
            :options="generationOptions"
            placeholder="ជ្រើសរើសជំនាន់"
            clearable
            label="ជំនាន់"
            :disabled="!formAcademicId"
          />
        </el-col>
        <el-col :span="8">
          <AppSelect
            v-model="form.term_id"
            :options="termOptions"
            placeholder="ជ្រើសរើសវគ្គ"
            clearable
            prop="term_id"
            label="វគ្គ"
            :disabled="!formGenerationId"
          />
        </el-col>
      </el-row>

      <!-- cascading select down to major -->
      <el-row :gutter="20">
        <el-col :span="6">
          <AppSelect
            v-model="formProgramID"
            :options="programmesOptions"
            placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
            clearable
            label="កម្មវិធីសិក្សា"
          />
        </el-col>
        <el-col :span="6">
          <AppSelect
            v-model="formFacultyID"
            :options="formFacultyOptions"
            placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
            clearable
            label="មហាវិទ្យាល័យ"
            :disabled="!formProgramID"
          />
        </el-col>
        <el-col :span="6">
          <AppSelect
            v-model="formDepartmentID"
            :options="formDepartmentOptions"
            placeholder="ជ្រើសរើសដេប៉ាតឺម៉ង់"
            clearable
            label="ដេប៉ាតឺម៉ង់"
            :disabled="!formFacultyID"
          />
        </el-col>
        <el-col :span="6">
          <AppSelect
            v-model="form.major_id"
            :options="formMajorOptions"
            placeholder="ជ្រើសរើសជំនាញ"
            clearable
            prop="major_id"
            label="ជំនាញ"
            :disabled="!formDepartmentID"
          />
        </el-col>
      </el-row>

      <!-- detail rows -->
      <div class="detail-section">
        <div class="detail-header">
          <el-text tag="b">ព័ត៌មានលម្អិត (ឆមាស / ឆ្នាំសិក្សា / វេន)</el-text>
          <AppButton type="default" icon="Plus" size="small" @click="addDetailRow">ថែមជួរ</AppButton>
        </div>

        <div
          v-for="(row, index) in form.class_curriclumn_details"
          :key="index"
          class="detail-row"
        >
          <div class="detail-row-header">
            <el-text tag="b" type="info">ជួរទី {{ index + 1 }}</el-text>
            <AppButton
              icon="Delete"
              circle
              size="small"
              type="default"
              plain
              :disabled="form.class_curriclumn_details.length <= 1"
              @click="removeDetailRow(index)"
            />
          </div>

          <el-row :gutter="20">
            <el-col :span="8">
              <AppSelect
                v-model="row.study_year_id"
                :options="studyyearOption"
                placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
                clearable
                label="ឆ្នាំសិក្សា"
              />
            </el-col>
            <el-col :span="8">
              <AppSelect
                v-model="row.semester_id"
                :options="semesterOptions"
                placeholder="ជ្រើសរើសឆមាស"
                clearable
                label="ឆមាស"
                :disabled="!formAcademicId"
              />
            </el-col>
            <el-col :span="8">
              <AppSelect
                v-model="row.academic_shift_id"
                :options="academicShiftOptions"
                placeholder="ជ្រើសរើសវេន"
                clearable
                label="វេន"
                :disabled="!formAcademicId"
              />
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="8">
              <AppInput
                v-model="row.midterm_date"
                type="date"
                placeholder="ថ្ងៃប្រឡងកម្ពស់"
                clearable
                label="ថ្ងៃប្រឡងកម្ពស់"
              />
            </el-col>
            <el-col :span="8">
              <AppInput
                v-model="row.final_date"
                type="date"
                placeholder="ថ្ងៃប្រឡងចុងក្រោយ"
                clearable
                label="ថ្ងៃប្រឡងចុងក្រោយ"
              />
            </el-col>
            <el-col :span="8">
              <AppSelect
                v-model="row.type_class"
                :options="typeClassOptions"
                placeholder="ជ្រើសរើសប្រភេទថ្នាក់"
                clearable
                label="ប្រភេទថ្នាក់"
              />
            </el-col>
          </el-row>
        </div>
      </div>
    </AppForm>
  </AppDialog>
</template>

<style scoped>
.detail-section {
  margin-top: 16px;
  border-top: 1px solid #ebeef5;
  padding-top: 16px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.detail-row {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 12px;
  background: #fafafa;
}

.detail-row-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
</style>