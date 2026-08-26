<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import {
  CreateClassCurriculum,
  GetClassCurriculumWithTeacherRate,
} from '../../services/classcurriculmn.service.js'
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
import TableCustom from '../../components/tables/TableCustom.vue'
import { CreateClassOffering } from '../../services/class_offering.service.js'
import { getSubjectByMajor } from '../../services/subject.service.js'
import { studentTermGet } from '../../services/studentterm.service.js'
import { createcourseregistration } from '../../services/course_registration.service.js'
import { getTeacherFilter } from '../../services/teacher.service.js'
import { createTeacherRate } from '../../services/teacher.service.js'
const notify = useNotification()
const classcurriculmn = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const submitting = ref(false)
const formRef = ref(null)
const dialogVisible = ref(false)
const isEditMode = ref(false)

const studentcolumns = [
  { prop: 'name_kh', label: 'ឈ្មោះខ្មែរ', minwidth: 120 },
  { prop: 'name_en', label: 'ឈ្មោះអង់គ្លេស', minwidth: 120 },
  { prop: 'date_of_birth', label: 'ថ្ងៃ-ខែ-ឆ្នាំកំណើត', minwidth: 120 },
  { prop: 'gender', slot: 'gender', label: 'ភេទ', minwidth: 120 },
  { prop: 'nationality', label: 'សញ្ជាតិ', minwidth: 120 },
  { prop: 'phone', label: 'លេខទូរសព្ទ', minwidth: 120 },
  { prop: 'occupation', label: 'មុខរបរ', minwidth: 120 },
]

const columns = [
  { prop: 'name', label: 'ឈ្មោះថ្នាក់', minwidth: 120 },
  { prop: 'major_name', slot: 'major_name', label: 'ជំនាញ', minwidth: 120 },
  {
    prop: 'major_duration_interval',
    slot: 'major_duration_interval',
    label: 'រយៈពេលសិក្សា',
    minwidth: 120,
  },
  { slot: 'programme_name', label: 'កម្រិត', minwidth: 120 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 120 },
  { prop: 'generation_name', label: 'ជំនាន់', minwidth: 120 },
  { prop: 'term_name', label: 'វគ្គ', minwidth: 120 },
  { prop: 'active', slot: 'active', label: 'ស្ថានភាព', minwidth: 120 },
]

const columndetails = [
  { prop: 'semester_name', label: 'ឆមាស', minwidth: 120 },
  { prop: 'study_year_id', label: 'ឆ្នាំទី', minwidth: 120 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 120 },
  { prop: 'academic_shift_name', label: 'វេនសិក្សា', minwidth: 120 },
  { prop: 'midterm_date', label: 'ថ្ងៃប្រឡង Midterm', minwidth: 120 },
  { prop: 'final_date', label: 'ថ្ងៃប្រឡង Final', minwidth: 120 },
  { prop: 'total_student', label: 'ប្រធានថ្នាក់', minwidth: 120 },
  { prop: 'type_class', slot: 'type_class', label: 'ប្រភេទថ្នាក់', minwidth: 120 },
]

const columnclass_offering = [
  { prop: 'subject_name', slot: 'subject_name', label: 'ឈ្មោះ', minminwidth: 120 },
  { prop: 'credit', slot: 'credit', label: 'ក្រេឌីត', minwidth: 90 },
  { prop: 'passing_score', slot: 'passing_score', label: 'ពិន្ទុជាប់', minwidth: 90 },
  { prop: 'total_hour', slot: 'total_hour', label: 'ម៉ោងសរុប', minwidth: 90 },
  { prop: 'total_hour', slot: 'total_hour', label: 'ម៉ោងនៅសល់', minwidth: 140 },
  { prop: 'status', label: 'ស្ថានភាព', minwidth: 90 },
  { prop: 'teacher_name',slot:'teacher_name', label: 'គ្រូបង្រៀន', minwidth: 170 },
  { prop: 'hourly_rate',slot:'hourly_rate', label: 'តម្លៃម៉ោងគ្រូ', minwidth: 90 },
  { prop: 'effective_date', label: 'ថ្ងៃខែមានប្រសិទ្ធភាព', minwidth: 90 },
]

async function fetchClassCurriculum() {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    const res = await GetClassCurriculumWithTeacherRate(params)
    classcurriculmn.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load class curriculums')
  }
}

/* ---------------- Class Offering dialog ---------------- */
const createclassofferingvisible = ref(false)
const submittingOffering = ref(false)
const subjectoptions = ref([])

function emptyOfferingRow() {
  return {
    subject_id: null,
    credit: 0,
    passing_score: 0,
    total_hour: 0,
    total_attendance_for_rexam: 0,
    total_attendance_for_relearn: 0,
    description: '',
  }
}

function defaultOfferingForm() {
  return {
    class_curriculum_detail_id: null,
    class_offering: [emptyOfferingRow()],
  }
}

const classofferingform = reactive(defaultOfferingForm())

function newClassOffering() {
  classofferingform.class_offering.push(emptyOfferingRow())
}

function removeClassOfferingRow(index) {
  if (classofferingform.class_offering.length <= 1) return
  classofferingform.class_offering.splice(index, 1)
}

// detailRow = the row from class_curriculum_detais (has the detail id)
// curriculumRow = the parent class curriculum row (has major_id)
async function openClassOfferingDialog(detailRow, curriculumRow) {
  Object.assign(classofferingform, defaultOfferingForm())
  classofferingform.class_curriculum_detail_id = detailRow.id
  subjectoptions.value = await loadSubjectOption(curriculumRow.major_id)
  createclassofferingvisible.value = true
}

function closeClassOfferingDialog() {
  createclassofferingvisible.value = false
  subjectoptions.value = []
  Object.assign(classofferingform, defaultOfferingForm())
}

async function submitClassOffering() {
  submittingOffering.value = true
  try {
    await CreateClassOffering(classofferingform)
    notify.success('បញ្ចូលមុខវិជ្ជាបានជោគជ័យ')
    closeClassOfferingDialog()
    fetchClassCurriculum()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to create class offering')
  } finally {
    submittingOffering.value = false
  }
}

/* ---------------- Class Curriculum top-level form ---------------- */


// class registration start
const createclassregistrationvisible = ref(false)
const selectedRows = ref([]);
const studentterms = ref([])
const selectedClassOffering = ref(null)

const classregistrationfilters = reactive({
  semester_id: null,
  study_year_id: null,
  term_id: null,
  major_id: null
})

const columnclass_registration = [
  { prop: 'code', label: 'លេខកូដ', minwidth: 90 },
  { prop: 'name', label: 'ឈ្មោះ', minwidth: 90 },
  { prop: 'gender', label: 'ភេទ', minwidth: 90 },
  { prop: 'phone', label: 'លេខទូរសព្ទ', minwidth: 90 },
  { prop: 'email', label: 'អុីម៉ែល', width: 200 },
  { prop: 'date_of_birth', label: 'ថ្ងៃ-ខែ-ឆ្នាំ', minwidth: 90 },
  { prop: 'place_of_birth', label: 'កន្លែងកំណេីត', minwidth: 90 },
  { prop: 'nationality', label: 'សញ្ជាតិ', minwidth: 90 },
  { prop: 'address', label: 'អាស័យដ្ឋាន', minwidth: 90 },
]

async function fetchTeacher(filters) {
  try {
    const params = {}

    if (filters.faculty_id) {
      params.faculty_id = filters.faculty_id
    }

    const res = await getTeacherFilter(params)

    return (res.data.data || []).map((f) => ({
      label: `${f.name} - (${f.gender})`,
      value: f.id
    }))
  } catch (e) {
    notify.error(
      e?.response?.data?.message ||
      e.message ||
      'Failed to load students'
    )

    return []
  }
}

const teacherrateform = reactive({
  teacher_id: null,
  class_offer_id: null,
  hourly_rate: null,
  effective_date: ''
})

async function openTeacherRateDialog(
  detailRow,
  curriculumRow,
  class_offering
) {
  classregistrationfilters.faculty_id = curriculumRow.faculty_id

  selectedClassOffering.value = class_offering
  selectedRows.value = []

  studentterms.value = await fetchTeacher(
    classregistrationfilters
  )

  createclassregistrationvisible.value = true
}

function closeTeacherRateDialog() {
  createclassregistrationvisible.value = false
  studentterms.value = []
}

async function submitteacherrate() {
  try {
    const payload = {
      teacher_id: teacherrateform.teacher_id,
      class_offer_id: selectedClassOffering.value.id,
      hourly_rate: teacherrateform.hourly_rate,
      effective_date: teacherrateform.effective_date
    }

    await createTeacherRate(payload)

    notify.success('បញ្ចូលគ្រូបានជោគជ័យ')

    selectedRows.value = []
    closeTeacherRateDialog()
    fetchClassCurriculum()
  } catch (e) {
    notify.error(
      e?.response?.data?.message ||
      e.message ||
      'Failed to create course registration'
    )
  }
}

function emptyDetailRow() {
  return {
    semester_id: null,
    study_year_id: null,
    academic_shift_id: null,
    midterm_date: '',
    final_date: '',
    type_class: '',
  }
}

function defaultForm() {
  return {
    name: '',
    major_id: null,
    term_id: null,
    class_curriclumn_details: [emptyDetailRow()],
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

async function loadSubjectOption(majorID) {
  if (!majorID) return []
  try {
    const res = await getSubjectByMajor(majorID)
    return (res.data.data || []).map((m) => ({ label: m.name, value: m.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load subjects')
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
    fetchClassCurriculum()
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
  fetchClassCurriculum()
  fetchProgrammeOptions()
  fetchAcademicOptions()
})
</script>

<template>
  <AppFilterBar :fields="[{ slot: 'create', span: 4 }]">
    <template #create>
      <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្នាក់ថ្មី</AppButton>
    </template>
  </AppFilterBar>

  <TableCustom expandable :data="classcurriculmn" :columns="columns" :total="total" v-model:current-page="page"
    v-model:page-size="pageSize" @page-change="fetchClassCurriculum">
    <template #programme_name="{ row }">
      <el-text tag="b" style="color: crimson">
        {{ row.programme_name }}
      </el-text>
    </template>
    <template #active="{ row }">
      <el-text>
        {{ row.active === true ? 'កំពុងសិក្សា' : 'បានបញ្ចប់' }}
      </el-text>
    </template>
    <template #major_name="{ row }">
      <el-text tag="b"> {{ row.major_code }} - {{ row.major_name }} </el-text>
    </template>
    <template #major_duration_interval="{ row }">
      <el-text tag="b"> {{ row.major_duration_period }} {{ row.major_duration_interval }} </el-text>
    </template>
    <!-- `row` here is the parent class curriculum row -->
    <template #expand="{ row }">
      <el-divider content-position="left">លំអិត</el-divider>
      <TableCustom expandable :data="row.class_curriculum_detais" :columns="columndetails" :show-pagination="false">
        <template #type_class="{ row }">
          <el-text>
            {{ row.type_class === 'onclass' ? 'រៀនថ្នាក់ផ្ទាល់' : 'រៀនOnline' }}
          </el-text>
        </template>
        <!-- rename inner scope's `row` to `detailRow` so it doesn't shadow the outer curriculum `row` -->
        <template #actions="{ row: detailRow }">
          <AppButton @click="openClassOfferingDialog(detailRow, row)"> ថែមមុខវិជ្ជា </AppButton>
        </template>
        <template #expand="{ row: detailRow }">
          <el-divider content-position="left">
            <el-text tag="b" style="color: darkcyan;">
              មុខវិជ្ជាត្រូវសិក្សា
            </el-text>
          </el-divider>
          <TableCustom expandable :data="detailRow.class_offering" :columns="columnclass_offering"
            :show-pagination="false" actions-width="200px">
            <template #hourly_rate="{row}">
              <el-text tag="b" style="color: crimson;">
                {{ row.hourly_rate }} $
              </el-text>
            </template>
            <template #teacher_name="{row}">
              <el-text>
                {{ row.teacher_name }} ({{ row.teacher_gender }})
              </el-text>
            </template>
            <template #credit="{ row }">
              <el-text tag="b" style="color: crimson;">
                {{ row.credit }}
              </el-text>
            </template>
            <template #passing_score="{ row }">
              <el-text tag="b" style="color: crimson;">
                {{ row.passing_score }}
              </el-text>
            </template>
            <template #total_hour="{ row }">
              <el-text tag="b" style="color: black;">
                {{ row.total_hour }} ម៉ោង
              </el-text>
            </template>
            <template #subject_name="{ row }">
              <el-text tag="b" style="color: black;">
                {{ row.subject_name }}
              </el-text>
            </template>
            <template #total_attendance_for_rexam="{ row }">
              <el-text tag="b" style="color: crimson;">
                {{ row.total_attendance_for_rexam }}
              </el-text>
            </template>
            <template #total_attendance_for_relearn="{ row }">
              <el-text tag="b" style="color: crimson;">
                {{ row.total_attendance_for_relearn }}
              </el-text>
            </template>
            <template #actions="{ row: class_offering }">
              <AppButton size="small" type="success"
                @click="openTeacherRateDialog(detailRow, row, class_offering)">
                បញ្ចូលគ្រូ
              </AppButton>
              <AppButton size="small" type="primary">
                កាលវិភាគ
              </AppButton>
            </template>
            <template #expand="{ row: offeringRow }">
              <el-divider content-position="left">
                <el-text tag="b" style="color: darkcyan;">សិស្សសរុបកំពុងសិក្សា {{ offeringRow.student.length }}
                  នាក់</el-text>
              </el-divider>
              <TableCustom expandable :data="offeringRow.student" :columns="studentcolumns" :show-pagination="false">
                <template #gender="{ row }">
                  <el-text>
                    {{ row.gender === 'Male' ? 'ប្រុស' : 'ស្រី' }}
                  </el-text>
                </template>

              </TableCustom>
            </template>function emptyDetailRow() {
  return {
    semester_id: null,
    study_year_id: null,
    academic_shift_id: null,
    midterm_date: '',
    final_date: '',
    type_class: '',
  }
}

          </TableCustom>
        </template>

      </TableCustom>

    </template>
  </TableCustom>

  <AppDialog v-if="dialogVisible" v-model:visible="dialogVisible"
    :title="isEditMode ? 'កែប្រែព័ត៌មានកម្មវិធីសិក្សា' : 'បង្កើតកម្មវិធីសិក្សាថ្មី'" :showDefaultFooter="false"
    width="60%" @close="closeDialog">
    <AppForm ref="formRef" :model="form" :rules="rules" :loading="submitting" :show-actions="true"
      @submit="handleSubmit" submitText="រក្សាទុក" resetText="ចាកចេញ">
      <el-row :gutter="20">
        <el-col :span="12">
          <AppInput v-model="form.name" placeholder="បញ្ចូលឈ្មោះកម្មវិធីសិក្សា" clearable prop="name"
            label="ឈ្មោះកម្មវិធីសិក្សា" />
        </el-col>
      </el-row>

      <!-- cascading select: Academic -> Generation -> Term -->
      <el-row :gutter="20">
        <el-col :span="8">
          <AppSelect v-model="formAcademicId" :options="academicOptions" placeholder="ជ្រើសរើសឆ្នាំសិក្សា" clearable
            label="ឆ្នាំសិក្សា" />
        </el-col>
        <el-col :span="8">
          <AppSelect v-model="formGenerationId" :options="generationOptions" placeholder="ជ្រើសរើសជំនាន់" clearable
            label="ជំនាន់" :disabled="!formAcademicId" />
        </el-col>
        <el-col :span="8">
          <AppSelect v-model="form.term_id" :options="termOptions" placeholder="ជ្រើសរើសវគ្គ" clearable prop="term_id"
            label="វគ្គ" :disabled="!formGenerationId" />
        </el-col>
      </el-row>

      <!-- cascading select down to major -->
      <el-row :gutter="20">
        <el-col :span="6">
          <AppSelect v-model="formProgramID" :options="programmesOptions" placeholder="ជ្រើសរើសកម្មវិធីសិក្សា" clearable
            label="កម្មវិធីសិក្សា" />
        </el-col>
        <el-col :span="6">
          <AppSelect v-model="formFacultyID" :options="formFacultyOptions" placeholder="ជ្រើសរើសមហាវិទ្យាល័យ" clearable
            label="មហាវិទ្យាល័យ" :disabled="!formProgramID" />
        </el-col>
        <el-col :span="6">
          <AppSelect v-model="formDepartmentID" :options="formDepartmentOptions" placeholder="ជ្រើសរើសដេប៉ាតឺម៉ង់"
            clearable label="ដេប៉ាតឺម៉ង់" :disabled="!formFacultyID" />
        </el-col>
        <el-col :span="6">
          <AppSelect v-model="form.major_id" :options="formMajorOptions" placeholder="ជ្រើសរើសជំនាញ" clearable
            prop="major_id" label="ជំនាញ" :disabled="!formDepartmentID" />
        </el-col>
      </el-row>

      <!-- detail rows -->
      <div class="detail-section">
        <div class="detail-header">
          <el-text tag="b">ព័ត៌មានលម្អិត (ឆមាស / ឆ្នាំសិក្សា / វេន)</el-text>
          <AppButton type="default" icon="Plus" size="small" @click="addDetailRow">ថែមជួរ</AppButton>
        </div>

        <div v-for="(row, index) in form.class_curriclumn_details" :key="index" class="detail-row">
          <div class="detail-row-header">
            <el-text tag="b" type="info">ឆមាសទី {{ index + 1 }}</el-text>
            <AppButton icon="Delete" circle size="small" type="default" plain
              :disabled="form.class_curriclumn_details.length <= 1" @click="removeDetailRow(index)" />
          </div>

          <el-row :gutter="20">
            <el-col :span="8">
              <AppSelect v-model="row.study_year_id" :options="studyyearOption" placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
                clearable label="ឆ្នាំសិក្សា" />
            </el-col>
            <el-col :span="8">
              <AppSelect v-model="row.semester_id" :options="semesterOptions" placeholder="ជ្រើសរើសឆមាស" clearable
                label="ឆមាស" :disabled="!formAcademicId" />
            </el-col>
            <el-col :span="8">
              <AppSelect v-model="row.academic_shift_id" :options="academicShiftOptions" placeholder="ជ្រើសរើសវេន"
                clearable label="វេន" :disabled="!formAcademicId" />
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="8">
              <AppInput v-model="row.midterm_date" type="date" placeholder="ថ្ងៃប្រឡង Midterm" clearable
                label="ថ្ងៃប្រឡង Midterm" />
            </el-col>
            <el-col :span="8">
              <AppInput v-model="row.final_date" type="date" placeholder="ថ្ងៃប្រឡង Final" clearable
                label="ថ្ងៃប្រឡង Final" />
            </el-col>
            <el-col :span="8">
              <AppSelect v-model="row.type_class" :options="typeClassOptions" placeholder="ជ្រើសរើសប្រភេទថ្នាក់"
                clearable label="ប្រភេទថ្នាក់" />
            </el-col>
          </el-row>
        </div>
      </div>
    </AppForm>
  </AppDialog>

  <!-- Class Offering dialog -->
  <AppDialog v-if="createclassofferingvisible" v-model:visible="createclassofferingvisible"
    title="បញ្ចូលមុខវិជ្ជាទៅថ្នាក់រៀន" :showDefaultFooter="false" width="720px" @close="closeClassOfferingDialog">
    <AppForm :model="classofferingform" :loading="submittingOffering" :show-actions="true" @submit="submitClassOffering"
      submitText="បញ្ជូល">
      <div class="section-header">
        <el-text tag="b">បញ្ចូលមុខវិជ្ជា</el-text>
        <AppButton type="warning" plain icon="Plus" size="small" @click="newClassOffering">
          បន្ថែម
        </AppButton>
      </div>

      <el-card v-for="(classf, index) in classofferingform.class_offering" :key="index" :gutter="16"
        style="margin-bottom: 12px">
        <template #header>
          <div class="section-header">
            <span>មុខវិជ្ជាទី -{{ index + 1 }}</span>
            <AppButton type="danger" icon="Delete" size="small" plain
              :disabled="classofferingform.class_offering.length <= 1" @click="removeClassOfferingRow(index)">
              លុប
            </AppButton>
          </div>
        </template>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect v-model="classf.subject_id" :options="subjectoptions" placeholder="មុខវិជ្ជា" label="មុខវិជ្ជា"
              clearable />
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="classf.credit" placeholder="បញ្ចូលក្រេឌីត" type="number" clearable
              label="បញ្ចូលក្រេឌីត" />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput v-model.number="classf.passing_score" placeholder="បញ្ចូលពន្ទុជាប់" type="number" clearable
              label="បញ្ចូលពន្ទុជាប់" />
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="classf.total_hour" placeholder="បញ្ចូលចំនួនម៉ោងត្រូវរៀន" type="number" clearable
              label="បញ្ចូលចំនួនម៉ោងត្រូវរៀន" />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput v-model.number="classf.total_attendance_for_relearn" placeholder="បញ្ចូលចំនួនអវត្តមានត្រូវរៀនសង"
              type="number" clearable label="បញ្ចូលចំនួនអវត្តមានត្រូវរៀនសង" />
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="classf.total_attendance_for_rexam" placeholder="បញ្ចូលចំនួនវត្តមានត្រូវប្រឡងសង"
              type="number" clearable label="បញ្ចូលចំនួនវត្តមានត្រូវប្រឡងសង" />
          </el-col>
        </el-row>

        <AppInput v-model="classf.description" placeholder="ផ្សេងៗ" clearable label="ផ្សេងៗ" />
      </el-card>
    </AppForm>
  </AppDialog>
  <AppDialog v-if="createclassregistrationvisible" v-model:visible="createclassregistrationvisible" title="បញ្ចូលសិស្ស"
    :showDefaultFooter="false" width="55%" @close="closeTeacherRateDialog">
    <AppForm :show-actions="true" @submit="submitteacherrate" submitText="រក្សាទុក" resetText="ចាកចេញ">

      <el-row :gutter="20">
        <el-col :span="8">
          <AppSelect v-model="teacherrateform.teacher_id" :options="studentterms" placeholder="រេីសគ្រូ" label="រេីសគ្រូ" clearable />
        </el-col>

        <el-col :span="8">
          <AppInput v-model.number="teacherrateform.hourly_rate" placeholder="តម្លៃម៉ោងគ្រូ" label="តម្លៃម៉ោងគ្រូ" clearable type="number" />
        </el-col>


        <el-col :span="8">
          <el-form-item label="ថ្ងៃមានប្រសិទ្ធភាព" prop="date_of_birth">
            <el-date-picker v-model="teacherrateform.effective_date" type="date" value-format="YYYY-MM-DD" placeholder="ថ្ងៃមានប្រសិទ្ធភាព"
              style="width: 100%" />
          </el-form-item>
        </el-col>

      </el-row>

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

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
</style>
