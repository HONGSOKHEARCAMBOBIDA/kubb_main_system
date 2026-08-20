<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { createStudent, getStudent, updateStudent } from '../../services/student.service'
import { getFeediscountGroup } from '../../services/feediscountgroup.service'
import { getProvince, getDistrict, getCommunce, getVillage } from '../../services/location.service'
import { getAcademicStream } from '../../services/academic_stream'
import { getDocumentType } from '../../services/document_type.service'
import { useNotification } from '../../composables/useNotification'
import { getAcademics } from '../../services/academic.service.js'
import { getGenerationByAcademic } from '../../services/generation.service.js'
import { getTermByGeneation } from '../../services/term.service'
import { getAcademicDegreeByAcademic } from '../../services/academic_degree.service'
import { getSchoolarshipGroup } from '../../services/schoolarship.service.js'
import { getSemesterByAcademic } from '../../services/semester.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppForm from '@/components/forms/AppForm.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppTabs from '../../components/tap/AppTabs.vue'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'

// Update state
const isEditMode = ref(false)
const editingId = ref(null)
const dialogVisible = ref(false)
const notify = useNotification()
const submitting = ref(false)
const formRef = ref(null)
const activeTab = ref('general')
const genderOption = [
  { label: 'ស្រី', value: 'Female' },
  { label: 'ប្រុស', value: 'Male' },
]

const nationalityOption = [
  {label: 'ខ្មែរ',value: 'ខ្មែរ'},
  {label: 'ជនជាតិភាគតិច',value:'ជនជាតិភាគតិច'},
  {label: 'ជនបរទេស',value:'ជនបរទេស'}
]

// List state
const students = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filters = reactive({ name: '' })
let searchTimer = null

// Detail

const detailVisible = ref(false)
const detailStudent = ref(null)

async function openDetail(row) {
  detailVisible.value = true
  detailStudent.value = row
}

function closeDetail() {
  detailVisible.value = false
  detailStudent.value = null
}

const columns = [
  { prop: 'code', label: 'លេខកូដ', minwidth: 120 },
  { prop: 'name_kh', label: 'ឈ្មោះខ្មែរ', minwidth: 200 },
  { prop: 'name_en', label: 'ឈ្មោះឡាតាំង', minwidth: 200 },
  { prop: 'gender', label: 'ភេទ', slot: 'gender', minwidth: 90 },
  { prop: 'date_of_birth', label: 'ថ្ងៃខែឆ្នាំកំណើត', minwidth: 130 },
  { prop: 'phone', label: 'លេខទូរសព្ទ', minwidth: 130 },
  { prop: 'group_name', label: 'ប្រភេទនិស្សិត', minwidth: 150 },
  { prop: 'academic_stream_name', label: 'សញ្ញាបត្រ', minwidth: 130 },
  { prop: 'status', label: 'ស្ថានភាព', minwidth: 350, slot: 'status' },
]

const columneducationdetail = [
  { prop: 'school_name', label: 'សាលារៀន', width: 120 },
  { prop: 'level', label: 'កម្រិត', width: 120 },
  { prop: 'start_date', label: 'ចាប់ផ្ដេីម', width: 120 },
  { prop: 'end_date', label: 'បញ្ជប់', width: 120 },
  { prop: 'cerificate_date', label: 'ទទួលសញ្ញាបត្រ', width: 120 },
  { prop: 'score', label: 'SCORE', width: 120 },
  { prop: 'gpa', label: 'GPA', width: 120 },
  { prop: 'grade', label: 'GRADE', width: 120 },
  { prop: 'villlage_name_kh', label: 'ភូមិ', width: 120 },
  { prop: 'communce_name', label: 'ឃុំ/សង្កាត់', width: 120 },
  { prop: 'distirct_name', label: 'ស្រុក/ខណ្ឌ', width: 120 },
  { prop: 'province_name', label: 'ខេត្ត/រាជធានី', width: 120 },
]

const columndocumentdetail = [
  { prop: 'document_type_name_kh', label: 'ប្រភេទឯកសារ', minwidth: 200 },
  { prop: 'required_qty', label: 'ចំនួនត្រូវការ', minwidth: 200 },
  { prop: 'received_qty', label: 'ចំនួនបានប្រគល់ជួន', minwidth: 200 },
  { prop: 'remark', label: 'កត់ចំណាំផ្សេងៗ', minwidth: 200 },
]

function resetForm() {
  form.group_id = null
  form.name_kh = ''
  form.name_en = ''
  form.date_of_birth = null
  form.gender = null
  form.nationality = ''
  form.phone = ''
  form.occupation = ''
  form.academic_stream_id = null
  form.father_name = ''
  form.father_english_name = ''
  form.father_age = ''
  form.father_phone_number = ''
  form.father_occupation = ''
  form.father_workplace = ''
  form.father_is_alive = null
  form.mother_name = ''
  form.mother_english_name = ''
  form.mother_age = ''
  form.mother_phone_number = ''
  form.mother_occupation = ''
  form.mother_workplace = ''
  form.mother_is_alive = null
  formRef.value?.resetFields()
  form.village_id = null
  form.student_educations = []
  form.student_documents = []
  eduLocations.splice(0, eduLocations.length)
  formProvinceID.value = null
  formDistrictID.value = null
  formCommunceID.value = null
  formDistrictOptions.value = []
  formCommunceOptions.value = []
  formVillageOptions.value = []
  form.admission = {
    term_id: null,
    academic_degree_id: null,
    date: '',
    state: '',
    description: '',
    referral_school: '',
  }
  admissionAcademicID.value = null
  admissionGenerationID.value = null
  admissionGenerationOptions.value = []
  admissionTermOptions.value = []
  admissionDegreeOptions.value = []
}

function openCreate() {
  isEditMode.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

async function openEdit(row) {
  isEditMode.value = true
  editingId.value = row.id
  resetForm()

  form.group_id = row.group_id
  form.name_kh = row.name_kh
  form.name_en = row.name_en
  form.date_of_birth = row.date_of_birth
  form.gender = row.gender
  form.nationality = row.nationality
  form.phone = row.phone
  form.village_id = row.village_id
  form.occupation = row.occupation
  form.academic_stream_id = row.academic_stream_id

  formProvinceID.value = row.province_id ?? null
  formDistrictID.value = row.district_id ?? null
  formCommunceID.value = row.communce_id ?? null
  if (formProvinceID.value) {
    formDistrictOptions.value = await loadDistrictOption(formProvinceID.value)
  }
  if (formDistrictID.value) {
    formCommunceOptions.value = await loadCommunceOption(formDistrictID.value)
  }
  if (formCommunceID.value) {
    formVillageOptions.value = await loadVillageOption(formCommunceID.value)
  }

  form.student_family = []

  const family = row.student_family || []

  for (const f of family) {
    form.student_family.push({
      father_name: f.father_name || '',
      father_english_name: f.father_english_name || '',
      father_age: f.father_age || null,
      father_is_alive: f.father_is_alive ?? true,
      father_phone_number: f.father_phone_number || '',
      father_occupation: f.father_occupation || '',
      father_workplace: f.father_workplace || '',

      mother_name: f.mother_name || '',
      mother_english_name: f.mother_english_name || '',
      mother_age: f.mother_age || null,
      mother_is_alive: f.mother_is_alive ?? true,
      mother_phone_number: f.mother_phone_number || '',
      mother_occupation: f.mother_occupation || '',
      mother_workplace: f.mother_workplace || '',
    })
  }

  form.student_educations = []
  eduLocations.splice(0, eduLocations.length)
  const educations = row.student_educations || []
  for (const e of educations) {
    form.student_educations.push({
      level: e.level || '',
      school_name: e.school_name || '',
      village_id: e.village_id ?? null,
      start_date: e.start_date || '',
      end_date: e.end_date || '',
      cerificate_date: e.cerificate_date || '',
      score: e.score || '',
      gpa: e.gpa || '',
      grade: e.grade || '',
    })

    const loc = {
      province_id: e.province_id ?? null,
      district_id: e.district_id ?? null,
      commune_id: e.communce_id ?? null,
      provinceOptions: formProvinceOptions.value,
      districtOptions: [],
      communeOptions: [],
      villageOptions: [],
    }
    if (loc.province_id) loc.districtOptions = await loadDistrictOption(loc.province_id)
    if (loc.district_id) loc.communeOptions = await loadCommunceOption(loc.district_id)
    if (loc.commune_id) loc.villageOptions = await loadVillageOption(loc.commune_id)
    eduLocations.push(loc)
  }
  // if (educations.length === 0) {
  //   newEducationRow()
  // }

  // --- Documents ---
  form.student_documents = []
  const documents = row.student_documents || []
  for (const d of documents) {
    form.student_documents.push({
      document_type_id: d.document_type_id ?? null,
      required_qty: d.required_qty ?? 1,
      received_qty: d.received_qty ?? 1,
      remark: d.remark || '',
    })
  }
  // if (documents.length === 0) {
  //   newDocumentRow()
  // }

  activeTab.value = 'general'
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  isEditMode.value = false
  editingId.value = null
}

async function fetchStudent() {
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.name) params.name = filters.name

    const res = await getStudent(params)
    students.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load students')
  }
}

const schoolarshipOptions = ref([])
const semesterOptions = ref([])

const groupOptions = ref([])
const academicStreamOptions = ref([])
const documentTypeOptions = ref([])

const academicOptions = ref([])
const generationOptions = ref([])
const termOptions = ref([])
const academicDegreeOptions = ref([])
const programmesOptions = ref([])

const admissionAcademicID = ref(null)
const admissionGenerationID = ref(null)
const admissionGenerationOptions = ref([])
const admissionTermOptions = ref([])
const admissionDegreeOptions = ref([])


const filteradmissionDegreeOptions = reactive({
  programme_id: null,
  faculty_id: null,
  department_id: null,
  major_id: null,
  academic_id: null,
})

const admissionStateOptions = [
  { label: 'កំពុងរង់ចាំ', value: 'created' },
  { label: 'ត្រូវបានទទួលយក', value: 'approved' },
  { label: 'បដិសេធ', value: 'rejected' },
]

const endrollmentFeeintervalOption = [
  { label: 'បង់1ខែម្ដង', value: 'monthly_fee' },
  { label: 'បង់3ខែម្ដង', value: 'quarterly_fee' },
  { label: 'បង់1ឆមាសម្ដង', value: 'semesterly_fee' },
  { label: 'បង់1ឆ្នាំម្ដង', value: 'yearly_fee' },
]

const studyyearOption = [
  { label: 'ឆ្នាំទី1', value: 1 },
  { label: 'ឆ្នាំទី2', value: 2 },
  { label: 'ឆ្នាំទី3', value: 3 },
  { label: 'ឆ្នាំទី4', value: 4 },
]

async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load programmes')
  }
}

async function fetchSchoolarship() {
  try {
    const res = await getSchoolarshipGroup()
    schoolarshipOptions.value = (res.data.data || []).map((a) => ({
      label:
        a.discount_type === 'percentage'
          ? `${a.name} ${a.discount_percentage}%`
          : `${a.name} ${a.discount_amount}$`,
      value: a.id,
    }))
  } catch (e) {
    console.error(e)
  }
}

async function fetchAcademicOption() {
  try {
    const res = await getAcademics()
    academicOptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }))
  } catch (e) {}
}

async function loadAcademicdegree() {
  try {
    const params = {}
    if (filteradmissionDegreeOptions.programme_id) params.programme_id = filteradmissionDegreeOptions.programme_id
    if (filteradmissionDegreeOptions.academic_id) params.academic_id = filteradmissionDegreeOptions.academic_id
    const res = await getAcademicDegreeByAcademic(params)
    admissionDegreeOptions.value = (res.data.data || []).map((g) => ({
      label: `${g.major_code} - ${g.programme_name} មួយឆ្នាំ ${g.yearly_fee}$`,
      value: g.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load districts')
    return []
  }
}

async function loadSemester(academicID) {
  if (!academicID) return []
  try {
    const res = await getSemesterByAcademic(academicID)
    return (res.data.data || []).map((g) => ({
      label: g.name,
      value: g.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load districts')
    return []
  }
}

async function loadGeneration(academicID) {
  if (!academicID) return []
  try {
    const res = await getGenerationByAcademic(academicID)
    return (res.data.data || []).map((g) => ({ label: g.name, value: g.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load districts')
    return []
  }
}

async function loadTerm(geneationID) {
  if (!geneationID) return []
  try {
    const res = await getTermByGeneation(geneationID)
    return (res.data.data || []).map((g) => ({ label: g.name, value: g.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load districts')
    return []
  }
}

async function onAdmissionAcademicChange() {
  admissionGenerationID.value = null
  form.admission.term_id = null
  admissionTermOptions.value = []
  form.admission.academic_degree_id = null
  academicDegreeOptions.value = []
 // admissionDegreeOptions.value = await loadAcademicdegree(admissionAcademicID.value)
  admissionGenerationOptions.value = await loadGeneration(admissionAcademicID.value)
  semesterOptions.value = await loadSemester(admissionAcademicID.value)
  filteradmissionDegreeOptions.academic_id = admissionAcademicID.value
}

async function onAdmissionGenerationChange() {
  form.admission.term_id = null
  admissionTermOptions.value = await loadTerm(admissionGenerationID.value)
}

async function fetchGroupOptions() {
  try {
    const res = await getFeediscountGroup()
    groupOptions.value = (res.data.data || []).map((g) => ({ 
    label: g.discount_type === 'percentage' ? `${g.name} បញ្ចុះតម្លៃ ${g.discount_percentage}%` : `${g.name} បញ្ចុះតម្លៃ ${g.discount_amount}$`, 
    value: g.id 
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load groups')
  }
}

async function fetchAcademicStreamOptions() {
  try {
    const res = await getAcademicStream()
    academicStreamOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academic streams')
  }
}

async function fetchDocumentTypeOptions() {
  try {
    const res = await getDocumentType()
    documentTypeOptions.value = (res.data.data || []).map((d) => ({
      label: d.name_kh,
      value: d.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load document types')
  }
}

async function loadDistrictOption(provinceID) {
  if (!provinceID) return []
  try {
    const res = await getDistrict(provinceID)
    return (res.data.data || []).map((f) => ({ label: f.name_kh, value: f.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load districts')
    return []
  }
}

async function loadCommunceOption(districtID) {
  if (!districtID) return []
  try {
    const res = await getCommunce(districtID)
    return (res.data.data || []).map((f) => ({ label: f.name_kh, value: f.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load communes')
    return []
  }
}

async function loadVillageOption(communceID) {
  if (!communceID) return []
  try {
    const res = await getVillage(communceID)
    return (res.data.data || []).map((f) => ({ label: f.name_kh, value: f.id }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load villages')
    return []
  }
}

const formProvinceID = ref(null)
const formDistrictID = ref(null)
const formCommunceID = ref(null)

const formProvinceOptions = ref([])
const formDistrictOptions = ref([])
const formCommunceOptions = ref([])
const formVillageOptions = ref([])

async function fetchProvinceOption() {
  try {
    const res = await getProvince()
    formProvinceOptions.value = (res.data.data || []).map((a) => ({
      label: a.name_kh,
      value: a.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load provinces')
  }
}

async function onProvinceChange() {
  formDistrictID.value = null
  formCommunceID.value = null
  form.village_id = null
  formCommunceOptions.value = []
  formVillageOptions.value = []
  formDistrictOptions.value = await loadDistrictOption(formProvinceID.value)
}

async function onDistrictChange() {
  formCommunceID.value = null
  form.village_id = null
  formVillageOptions.value = []
  formCommunceOptions.value = await loadCommunceOption(formDistrictID.value)
}

async function onCommunceChange() {
  form.village_id = null
  formVillageOptions.value = await loadVillageOption(formCommunceID.value)
}

// ---------------------------------------------------------------------------
// Main form model (mirrors StudentRequestCreate)
// ---------------------------------------------------------------------------
const form = reactive({
  group_id: null,

  name_kh: '',
  name_en: '',
  date_of_birth: '',
  gender: '',
  nationality: '',
  phone: '',
  village_id: null,
  occupation: '',
  academic_stream_id: null,

  student_educations: [],
  student_documents: [],

  student_family: [],

  admission: {
    term_id: null,
    academic_degree_id: null,
    date: '',
    state: '',
    description: '',
    referral_school: '',
  },
  enrollment: {
    scholarship_id: null,
    fee_interval: '',
  },
  student_term: {
    semester_id: null,
    study_year_id: null,
  },
})

const rules = {
  group_id: [{ required: true, message: 'Please select group', trigger: 'change' }],
  name_kh: [
    { required: true, message: 'Please enter Khmer name', trigger: 'blur' },
    {
      min: 2,
      max: 150,
      message: 'Khmer name must be between 2 and 150 characters',
      trigger: 'blur',
    },
  ],
  name_en: [
    { required: true, message: 'Please enter English name', trigger: 'blur' },
    {
      min: 2,
      max: 150,
      message: 'English name must be between 2 and 150 characters',
      trigger: 'blur',
    },
  ],
  date_of_birth: [{ required: true, message: 'Please select date of birth', trigger: 'change' }],
  gender: [{ required: true, message: 'Please select gender', trigger: 'change' }],
  nationality: [
    { required: true, message: 'Please enter nationality', trigger: 'blur' },
    {
      min: 2,
      max: 150,
      message: 'Nationality must be between 2 and 150 characters',
      trigger: 'blur',
    },
  ],
  phone: [
    { required: true, message: 'Please enter phone number', trigger: 'blur' },
    {
      min: 2,
      max: 150,
      message: 'Phone number must be between 2 and 150 characters',
      trigger: 'blur',
    },
  ],
  village_id: [{ required: true, message: 'Please select village', trigger: 'change' }],
  academic_stream_id: [
    { required: true, message: 'Please select academic stream', trigger: 'change' },
  ],
}

// ---------------------------------------------------------------------------
// Student educations (repeatable, each with its own address cascade)
// ---------------------------------------------------------------------------
// eduLocations[i] holds the cascade state for student_educations[i]
const eduLocations = reactive([])

function newEducationRow() {
  form.student_educations.push({
    level: '',
    school_name: '',
    village_id: null,
    start_date: '',
    end_date: '',
    cerificate_date: '',
    score: '',
    gpa: '',
    grade: '',
  })
  eduLocations.push({
    province_id: null,
    district_id: null,
    commune_id: null,
    provinceOptions: formProvinceOptions.value, // shares the same province list
    districtOptions: [],
    communeOptions: [],
    villageOptions: [],
  })
}

function removeEducationRow(index) {
  form.student_educations.splice(index, 1)
  eduLocations.splice(index, 1)
}

async function onEduProvinceChange(index) {
  const loc = eduLocations[index]
  loc.district_id = null
  loc.commune_id = null
  form.student_educations[index].village_id = null
  loc.communeOptions = []
  loc.villageOptions = []
  loc.districtOptions = await loadDistrictOption(loc.province_id)
}

async function onEduDistrictChange(index) {
  const loc = eduLocations[index]
  loc.commune_id = null
  form.student_educations[index].village_id = null
  loc.villageOptions = []
  loc.communeOptions = await loadCommunceOption(loc.district_id)
}

async function onEduCommunceChange(index) {
  const loc = eduLocations[index]
  form.student_educations[index].village_id = null
  loc.villageOptions = await loadVillageOption(loc.commune_id)
}

// ---------------------------------------------------------------------------
// Student documents (repeatable)
// ---------------------------------------------------------------------------
function newDocumentRow() {
  form.student_documents.push({
    document_type_id: null,
    required_qty: 1,
    received_qty: 1,
    remark: '',
  })
}

function newFamilyRow() {
  form.student_family.push({
    father_name: '',
    father_english_name: '',
    father_age: null,
    father_is_alive: true,
    father_phone_number: '',
    father_occupation: '',
    father_workplace: '',

    mother_name: '',
    mother_english_name: '',
    mother_age: null,
    mother_is_alive: true,
    mother_phone_number: '',
    mother_occupation: '',
    mother_workplace: '',
  })
}

function removeDocumentRow(index) {
  form.student_documents.splice(index, 1)
}

function removeFamilyRow(index) {
  form.student_family.splice(index, 1)
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditMode.value) {
      await updateStudent(editingId.value, form)
      //notify.success('Student updated successfully')
    } else {
      await createStudent(form)
      //notify.success('Student created successfully')
    }
    dialogVisible.value = false
    resetForm()
    isEditMode.value = false
    editingId.value = null
    fetchStudent()
  } catch (e) {
    notify.error(
      e?.response?.data?.message ||
        e.message ||
        (isEditMode.value ? 'Failed to update student' : 'Failed to create student')
    )
  } finally {
    submitting.value = false
  }
}

watch(
  () => filteradmissionDegreeOptions.programme_id,
  () => {
    loadAcademicdegree()
  }
)

watch(
  () => filteradmissionDegreeOptions.academic_id,
  () => {
    loadAcademicdegree()
  }
)

watch(
  () => filters.name,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchStudent()
    }, 400) // debounce so we don't hit the API on every keystroke
  }
)

onMounted(() => {
  fetchProvinceOption()
  fetchGroupOptions()
  fetchAcademicStreamOptions()
  fetchDocumentTypeOptions()
  // newEducationRow()
  // newDocumentRow()
  fetchStudent()
  fetchAcademicOption()
  fetchSchoolarship()
  fetchProgrammeOptions()
})
</script>

<template>
  <div class="student-page">
    <AppFilterBar
      :fields="[
        { slot: 'name', span: 8 },
        { slot: 'create', span: 4 },
      ]"
    >
      <template #name>
        <AppInput v-model="filters.name" placeholder="ស្វែងរកតាមឈ្មោះខ្មែរ" clearable />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតនិស្សិតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      :data="students"
      :columns="columns"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchStudent"
    >
      <template #gender="{ row }">
        {{ row.gender === 'Male' ? 'ប្រុស' : 'ស្រី' }}
      </template>

      <template #status="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.status }}
        </el-text>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="មើលព័ត៌មានលម្អិត" placement="top">
          <AppButton
            icon="View"
            circle
            size="small"
            type="success"
            plain
            @click="openDetail(row)"
          />
        </el-tooltip>
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="warning" plain @click="openEdit(row)" />
        </el-tooltip>
      </template>
    </TableCustom>
    <AppDialog
      v-if="dialogVisible"
      v-model:visible="dialogVisible"
      :title="isEditMode ? 'កែប្រែព័ត៌មាននិស្សិត' : 'បង្កេីតនិស្សិតថ្មី'"
      :showDefaultFooter="false"
      width="58%"
      @close="closeDialog"
    >
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :show-actions="true"
        @submit="handleSubmit"
        submitText="រក្សាទុក"
      >
        <AppTabs
          v-model="activeTap"
          :tabs="[
            { name: 'general', label: 'ព័ត៏មានទូទៅ' },
            { name: 'admission', label: 'ដាក់ពាក្យ' },
            { name: 'education', label: 'ការសិក្សា' },
            { name: 'document', label: 'ឯកសារប្រគល់ជួន' },
            // { name: 'fother', label: 'ព័ត៏មានឪពុក' },
            // { name: 'mother', label: 'ព័ត៏មានម្ដាយ' },
            { name: 'family', label: 'ព័ត៍មានគ្រូសារ' },
          ]"
          tab-position="top"
          stretch="true"
          :lazy="true"
        >
          <template #admission>
            <el-card class="section-card" shadow="never">
              <template #header>ដាក់ពាក្យចូលរៀន</template>
              <el-row :gutter="20">
                <el-col :span="8">
                  <AppSelect
                    v-model="admissionAcademicID"
                    :options="academicOptions"
                    placeholder="រេីសឆ្នាំសិក្សា"
                    label="ឆ្នាំសិក្សា"
                    clearable
                    @change="onAdmissionAcademicChange"
                  />
                </el-col>
                <el-col :span="8">
                  <AppSelect
                    v-model="filteradmissionDegreeOptions.programme_id"
                    :options="programmesOptions"
                    placeholder="កម្រិត"
                    label="កម្រិត"
                    clearable
                    :disabled="!admissionAcademicID"
                  />
                </el-col>
                <el-col :span="8">
                  <AppSelect
                    v-model="form.admission.academic_degree_id"
                    :options="admissionDegreeOptions"
                    placeholder="សុំចូលរៀនជំនាញ"
                    label="សុំចូលរៀនជំនាញ"
                    clearable
                    
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="admissionGenerationID"
                    :options="admissionGenerationOptions"
                    placeholder="រេីសជំនាន់"
                    label="ជំនាន់"
                    clearable
                    :disabled="!admissionAcademicID"
                    @change="onAdmissionGenerationChange"
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="form.admission.term_id"
                    :options="admissionTermOptions"
                    placeholder="រេីសវគ្គ"
                    label="រេីសវគ្គ"
                    clearable
                    :disabled="!admissionGenerationID"
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="ថ្ងៃដាក់ពាក្យ">
                    <el-date-picker
                      v-model="form.admission.date"
                      type="date"
                      value-format="YYYY-MM-DD"
                      placeholder="ថ្ងៃដាក់ពាក្យ"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="form.admission.state"
                    :options="admissionStateOptions"
                    placeholder="ស្ថានភាពដាក់ពាក្យ"
                    label="ស្ថានភាព"
                    clearable
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="form.enrollment.scholarship_id"
                    :options="schoolarshipOptions"
                    placeholder="អាហារូបករណ៍"
                    label="អាហារូបករណ៍"
                    clearable
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="form.enrollment.fee_interval"
                    :options="endrollmentFeeintervalOption"
                    placeholder="សុំបង់ប្រាក់ជា"
                    label="សុំបង់ប្រាក់ជា"
                    clearable
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppInput
                    v-model="form.admission.referral_school"
                    label="មកពីសាលា(មិនចាំចាច់)"
                    placeholder="មកពីសាលា"
                    clearable
                  />
                </el-col>
                <el-col :span="12">
                  <AppInput
                    v-model="form.admission.description"
                    label="សេចក្តីលម្អិត(មិនចាំចាច់)"
                    placeholder="សេចក្តីលម្អិត"
                    clearable
                  />
                </el-col>
              </el-row>
              <el-divider content-position="left">សុំចូលរៀន</el-divider>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="form.student_term.study_year_id"
                    :options="studyyearOption"
                    placeholder="ឆ្នាំទី"
                    label="ឆ្នាំទី"
                    clearable
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="form.student_term.semester_id"
                    :options="semesterOptions"
                    :disabled="!admissionAcademicID"
                    placeholder="ឆមាសទី"
                    label="ឆមាសទី"
                    clearable
                  />
                </el-col>
              </el-row>
            </el-card>
          </template>
          <template #general>
            <el-card class="section-card" shadow="never">
              <template #header>ព័ត៏មានទូទៅ</template>
              <el-row :gutter="20">
                <el-col :span="12">
                      <AppInput
                    v-model="form.name_kh"
                    prop="name_kh"
                    label="ឈ្មោះខែ្មរ"
                    placeholder="ឈ្មោះខែ្មរ"
                    clearable
                  />
              
                </el-col>
                <el-col :span="12">
                   <AppInput
                    v-model="form.name_en"
                    prop="name_en"
                    label="ឈ្មោះឡាតាំង"
                    placeholder="ឈ្មោះឡាតាំង"
                    clearable
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                      <AppSelect
                    v-model="form.group_id"
                    :options="groupOptions"
                    placeholder="ប្រភេទនិស្សិត"
                    prop="group_id"
                    label="ប្រភេទនិស្សិត"
                    clearable
                  />
                </el-col>
                <el-col :span="12">
                  <el-form-item label="ថ្ងៃខែឆ្នាំកំណេីត" prop="date_of_birth">
                    <el-date-picker
                      v-model="form.date_of_birth"
                      type="date"
                      value-format="YYYY-MM-DD"
                      placeholder="ថ្ងៃខែឆ្នាំកំណេីត"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    label="ភេទ"
                    v-model="form.gender"
                    :options="genderOption"
                    placeholder="រេីសភេទ"
                    clearable
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    label="សញ្ជាតិ"
                    v-model="form.nationality"
                    :options="nationalityOption"
                    placeholder="សញ្ជាតិ"
                    clearable
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <AppInput
                    v-model="form.phone"
                    prop="phone"
                    label="លេខទូរសព្ទ"
                    placeholder="លេខទូរសព្ទ"
                    type="number"
                    clearable
                  />
                </el-col>
                <el-col :span="8">
                  <AppInput
                    v-model="form.occupation"
                    prop="occupation"
                    label="មុខរបរ"
                    placeholder="មុខរបរ"
                    clearable
                  />
                </el-col>
                <el-col :span="8">
                  <AppSelect
                    v-model="form.academic_stream_id"
                    :options="academicStreamOptions"
                    placeholder="សញ្ញាបត្រ"
                    clearable
                    label="សញ្ញាបត្រ"
                    prop="academic_stream_id"
                  />
                </el-col>
              </el-row>
              <el-divider content-position="left">ទីកន្លែងរស់នៅ</el-divider>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="formProvinceID"
                    :options="formProvinceOptions"
                    placeholder="រេីសខេត្ត"
                    label="រេីសខេត្ត"
                    clearable
                    @change="onProvinceChange"
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="formDistrictID"
                    :options="formDistrictOptions"
                    placeholder="រេីសស្រុក"
                    :disabled="!formProvinceID"
                    clearable
                    label="រេីសស្រុក"
                    @change="onDistrictChange"
                  />
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <AppSelect
                    v-model="formCommunceID"
                    :options="formCommunceOptions"
                    placeholder="រេីសឃុំ"
                    :disabled="!formDistrictID"
                    clearable
                    label="រេីសឃុំ"
                    @change="onCommunceChange"
                  />
                </el-col>
                <el-col :span="12">
                  <AppSelect
                    v-model="form.village_id"
                    :options="formVillageOptions"
                    placeholder="រេីសភូមិ"
                    label="រេីសភូមិ"
                    :disabled="!formCommunceID"
                    clearable
                  />
                </el-col>
              </el-row>
            </el-card>
          </template>

          <template #education>
            <el-card class="section-card" shadow="never">
              <template #header>
                <div class="section-header">
                  <span>ការសិក្សា</span>
                  <AppButton type="warning" plain icon="Plus" size="small" @click="newEducationRow"
                    >បន្ថែមការសិក្សា</AppButton
                  >
                </div>
              </template>

              <el-card
                v-for="(edu, index) in form.student_educations"
                :key="index"
                class="sub-card"
                shadow="never"
              >
                <template #header>
                  <div class="section-header">
                    <span>ការសិក្សា -{{ index + 1 }}</span>
                    <AppButton
                      type="danger"
                      icon="Delete"
                      size="small"
                      plain
                      :disabled="form.student_educations.length === 1"
                      @click="removeEducationRow(index)"
                    >
                      លុប
                    </AppButton>
                  </div>
                </template>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppInput
                      v-model="edu.school_name"
                      label="ឈ្មោះសាលារៀន"
                      placeholder="ឈ្មោះសាលារៀន"
                      clearable
                    />
                  </el-col>
                  <el-col :span="12">
                    <AppInput
                      v-model="edu.level"
                      label="កម្រិតសិក្សា"
                      placeholder="កម្រិតសិក្សា"
                      clearable
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="8">
                    <AppInput v-model="edu.gpa" label="GPA" placeholder="GPA" clearable />
                  </el-col>
                  <el-col :span="8">
                    <AppInput v-model="edu.score" label="ពន្ទុ" placeholder="ពន្ទុ" clearable />
                  </el-col>
                  <el-col :span="8">
                    <AppInput v-model="edu.grade" label="Grade" placeholder="Grade" clearable />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="8">
                    <el-form-item label="ថ្ងៃខែឆ្នាំចាប់ផ្ដេីមសិក្សា">
                      <el-date-picker
                        v-model="edu.start_date"
                        type="date"
                        value-format="YYYY-MM-DD"
                        placeholder="ថ្ងៃខែឆ្នាំចាប់ផ្ដេីមសិក្សា"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="ថ្ងៃខែឆ្នាំបញ្ចប់សិក្សា">
                      <el-date-picker
                        v-model="edu.end_date"
                        type="date"
                        value-format="YYYY-MM-DD"
                        placeholder="ថ្ងៃខែឆ្នាំបញ្ចប់សិក្សា"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>

                  <el-col :span="8">
                    <el-form-item label="ថ្ងៃចេញសញ្ញាបត្រ">
                      <el-date-picker
                        v-model="edu.cerificate_date"
                        type="date"
                        value-format="YYYY-MM-DD"
                        placeholder="ថ្ងៃចេញសញ្ញាបត្រ"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-divider content-position="left">ទីកន្លែងសាលារៀន</el-divider>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppSelect
                      v-model="eduLocations[index].province_id"
                      :options="formProvinceOptions"
                      placeholder="រេីសខេត្ត"
                      label="រេីសខេត្ត"
                      clearable
                      @change="onEduProvinceChange(index)"
                    />
                  </el-col>
                  <el-col :span="12">
                    <AppSelect
                      v-model="eduLocations[index].district_id"
                      :options="eduLocations[index].districtOptions"
                      placeholder="រេីសស្រុក"
                      label="រេីសស្រុក"
                      :disabled="!eduLocations[index].province_id"
                      clearable
                      @change="onEduDistrictChange(index)"
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppSelect
                      v-model="eduLocations[index].commune_id"
                      :options="eduLocations[index].communeOptions"
                      placeholder="រេីសឃុំ"
                      label="រេីសឃុំ"
                      :disabled="!eduLocations[index].district_id"
                      clearable
                      @change="onEduCommunceChange(index)"
                    />
                  </el-col>
                  <el-col :span="12">
                    <AppSelect
                      v-model="edu.village_id"
                      :options="eduLocations[index].villageOptions"
                      placeholder="រេីសភូមិ"
                      label="រេីសភូមិ"
                      :disabled="!eduLocations[index].commune_id"
                      clearable
                    />
                  </el-col>
                </el-row>
              </el-card>
            </el-card>
          </template>

          <template #document>
            <el-card class="section-card" shadow="never">
              <template #header>
                <div class="section-header">
                  <span>ឯកសារនិស្សិតប្រគល់ជួន</span>
                  <AppButton type="warning" plain icon="Plus" size="small" @click="newDocumentRow"
                    >បន្ថែមឯកសារនិស្សិត</AppButton
                  >
                </div>
              </template>

              <el-card
                v-for="(doc, index) in form.student_documents"
                :key="index"
                :gutter="16"
                class="sub-card"
              >
                <template #header>
                  <div class="section-header">
                    <span>ឯកសារ -{{ index + 1 }}</span>
                    <AppButton
                      type="danger"
                      icon="Delete"
                      size="small"
                      plain
                      :disabled="form.student_documents.length === 1"
                      @click="removeDocumentRow(index)"
                    >
                      លុប
                    </AppButton>
                  </div>
                </template>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppSelect
                      v-model="doc.document_type_id"
                      :options="documentTypeOptions"
                      placeholder="ប្រភេទឯកសារ"
                      label="ប្រភេទឯកសារ"
                      clearable
                    />
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="ចំនួនត្រូវប្រគល់ជួន">
                      <el-input-number
                        v-model="doc.required_qty"
                        :min="1"
                        :max="20"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="ចំនួនទទួលបាន">
                      <el-input-number
                        v-model="doc.received_qty"
                        :min="1"
                        :max="20"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="សម្គាល់">
                      <el-input v-model="doc.remark" placeholder="សម្គាល់បេីមានបញ្ហា" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-card>
            </el-card>
          </template>
          <template #family>
            <el-card class="section-card" shadow="never">
              <template #header>
                <div class="section-header">
                  <span>ព័ត៏មានគ្រួសារ</span>
                  <AppButton type="warning" plain icon="Plus" size="small" @click="newFamilyRow"
                    >បន្ថែមព័ត៏មានគ្រួសារ</AppButton
                  >
                </div>
              </template>

              <el-card
                v-for="(fmy, index) in form.student_family"
                :key="index"
                :gutter="16"
                class="sub-card"
              >
                <template #header>
                  <div class="section-header">
                    <span>បន្ថែមព័ត៏មានគ្រួសារ -{{ index + 1 }}</span>
                    <AppButton
                      type="danger"
                      icon="Delete"
                      size="small"
                      plain
                      :disabled="form.student_family.length === 1"
                      @click="removeFamilyRow(index)"
                    >
                      លុប
                    </AppButton>
                  </div>
                </template>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.father_name"
                      label="ឈ្មោះខ្មែរឪពុក"
                      placeholder="ឈ្មោះខ្មែរឪពុក"
                      clearable
                    />
                  </el-col>
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.father_english_name"
                      label="ឈ្មោះឡាតាំងឪពុក"
                      placeholder="ឈ្មោះឡាតាំងឪពុក"
                      clearable
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="អាយុៈឪពុក">
                      <el-input-number
                        v-model="fmy.father_age"
                        :min="1"
                        :max="200"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.father_phone_number"
                      label="លេខទូរសព្ទឪពុក"
                      placeholder="លេខទូរសព្ទ"
                      type="number"
                      clearable
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="8">
                    <AppInput
                      v-model="fmy.father_occupation"
                      label="មុខរបរឪពុក"
                      placeholder="មុខរបរឪពុក"
                      clearable
                    />
                  </el-col>
                  <el-col :span="8">
                    <AppInput
                      v-model="fmy.father_workplace"
                      label="កន្លែងធ្វេីការឪពុក"
                      placeholder="កន្លែងធ្វេីការឪពុក"
                      clearable
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="នៅមានជីវិត">
                      <el-switch v-model="fmy.father_is_alive" />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.mother_name"
                      label="ឈ្មោះខ្មែរម្ដាយ"
                      placeholder="ឈ្មោះខ្មែរម្ដាយ"
                      clearable
                    />
                  </el-col>
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.mother_english_name"
                      label="ឈ្មោះឡាតាំងម្ដាយ"
                      placeholder="ឈ្មោះឡាតាំងម្ដាយ"
                      clearable
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="អាយុៈម្ដាយ">
                      <el-input-number
                        v-model="fmy.mother_age"
                        :min="1"
                        :max="200"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <AppInput
                      v-model="fmy.mother_phone_number"
                      label="លេខទូរសព្ទម្ដាយ"
                      placeholder="លេខទូរសព្ទម្ដាយ"
                      type="number"
                      clearable
                    />
                  </el-col>
                </el-row>
                <el-row :gutter="20">
                  <el-col :span="8">
                    <AppInput
                      v-model="fmy.mother_occupation"
                      label="មុខរបរម្ដាយ"
                      placeholder="មុខរបរម្ដាយ"
                      clearable
                    />
                  </el-col>
                  <el-col :span="8">
                    <AppInput
                      v-model="fmy.mother_workplace"
                      label="កន្លែងធ្វេីការម្ដាយ"
                      placeholder="កន្លែងធ្វេីការម្ដាយ"
                      clearable
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="នៅមានជីវិត">
                      <el-switch v-model="fmy.mother_is_alive" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-card>
            </el-card>
          </template>
        </AppTabs>
      </AppForm>
    </AppDialog>
    <AppDialog
      v-if="detailVisible"
      v-model:visible="detailVisible"
      :showDefaultFooter="false"
      width="75%"
    >
      <div>
        <div v-if="detailStudent">
          <div class="a4-header">
            <h2>ព័ត៌មានលម្អិតនិស្សិត</h2>
            <el-text tag="b" style="color: black">លេខកូដ: {{ detailStudent.code }}</el-text>
          </div>

          <section class="a4-section">
            <el-text tag="b" style="color: darkcyan" class="pb-5">ព័ត៌មានទូទៅ </el-text>
            <div class="a4-grid">
              <div><span class="a4-label">ឈ្មោះខ្មែរ</span> {{ detailStudent.name_kh }}</div>
              <div><span class="a4-label">ឈ្មោះឡាតាំង</span> {{ detailStudent.name_en }}</div>
              <div>
                <span class="a4-label">ភេទ</span>
                {{ detailStudent.gender === 'Male' ? 'ប្រុស' : 'ស្រី' }}
              </div>
              <div>
                <span class="a4-label">ថ្ងៃខែឆ្នាំកំណើត</span> {{ detailStudent.date_of_birth }}
              </div>
              <div><span class="a4-label">សញ្ជាតិ</span> {{ detailStudent.nationality }}</div>
              <div><span class="a4-label">លេខទូរសព្ទ</span> {{ detailStudent.phone }}</div>
              <div><span class="a4-label">មុខរបរ</span> {{ detailStudent.occupation || '-' }}</div>
              <div>
                <span class="a4-label">សញ្ញាបត្រ</span>
                {{ detailStudent.academic_stream_name || '-' }}
              </div>
              <div>
                <span class="a4-label">ស្ថិតក្នុងប្រភេទ</span>
                {{ detailStudent.group_name || '-' }}, បញ្ចុះតម្លៃ
                <el-text tag="b" style="color: crimson">
                  {{
                    detailStudent.discount_type === 'percentage'
                      ? `${detailStudent.discount_percentage}%`
                      : `${detailStudent.discount_amount}$`
                  }}
                </el-text>
              </div>
              <div><span class="a4-label">ស្ថានភាព</span> {{ detailStudent.status }}</div>
            </div>
            <div class="a4-address">
              <span class="a4-label">អាសយដ្ឋាន</span>
              ភូមិ{{ detailStudent.villlage_name_kh }}, ឃុំ{{ detailStudent.communce_name }},
              ស្រុក{{ detailStudent.distirct_name }}, ខេត្ត{{ detailStudent.province_name }}
            </div>
          </section>

          <section class="a4-section">
            <el-text tag="b" style="color: darkcyan" class="pb-5"> ព័ត៌មានឪពុកម្ដាយ </el-text>
            <div v-for="fam in detailStudent.student_family" :key="fam.id" class="a4-grid">
              <div><span class="a4-label">ឈ្មោះឪពុក</span> {{ fam.father_name || '-' }}</div>
              <div><span class="a4-label">អាយុឪពុក</span> {{ fam.father_age || '-' }}</div>
              <div><span class="a4-label">មុខរបរឪពុក</span> {{ fam.father_occupation || '-' }}</div>
              <div>
                <span class="a4-label">ទូរសព្ទឪពុក</span> {{ fam.father_phone_number || '-' }}
              </div>
              <div><span class="a4-label">ឈ្មោះម្ដាយ</span> {{ fam.mother_name || '-' }}</div>
              <div><span class="a4-label">អាយុម្ដាយ</span> {{ fam.mother_age || '-' }}</div>
              <div>
                <span class="a4-label">មុខរបរម្ដាយ</span> {{ fam.mother_occupation || '-' }}
              </div>
              <div>
                <span class="a4-label">ទូរសព្ទម្ដាយ</span> {{ fam.mother_phone_number || '-' }}
              </div>
            </div>
            <el-empty
              v-if="!detailStudent.student_family?.length"
              description="គ្មានទិន្នន័យ"
              :image-size="40"
            />
          </section>

          <section class="a4-section">
            <el-text tag="b" style="color: darkcyan" class="pb-5">ការសិក្សា </el-text>

            <TableCustom
              :data="detailStudent.student_educations"
              :columns="columneducationdetail"
              :show-pagination="false"
            >
            </TableCustom>
          </section>

          <section class="a4-section">
            <el-text tag="b" style="color: darkcyan">ឯកសារប្រគល់ជូន </el-text>
            <TableCustom
              :data="detailStudent.student_documents"
              :columns="columndocumentdetail"
              :show-pagination="false"
            >
            </TableCustom>
          </section>
        </div>
      </div>
    </AppDialog>
  </div>
</template>

<style scoped>
.student-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.section-card {
  margin-bottom: 16px;
}
.section-card :deep(.el-card__header) {
  font-weight: 600;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.sub-card {
  margin-bottom: 12px;
  background: var(--el-fill-color-light, #f7f8fa);
}
.doc-row {
  align-items: flex-end;
}
.doc-remove {
  display: flex;
  justify-content: flex-end;
  padding-bottom: 20px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 12px 0 24px;
}
.a4-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}
.a4-page {
  width: 210mm;
  min-height: 297mm;
  margin: 0 auto;
  padding: 15mm;
  box-sizing: border-box;
}
.a4-header {
  text-align: center;
  border-bottom: 2px solid #333;
  padding-bottom: 12px;
  margin-bottom: 16px;
}
.a4-header h2 {
  margin: 0;
}
.a4-code {
  color: #888;
  font-size: 13px;
  margin-top: 4px;
}
.a4-section {
  margin-bottom: 20px;
}
.a4-section h3 {
  font-size: 15px;
  margin-bottom: 10px;
}
.a4-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  font-size: 13px;
}
.a4-label {
  color: #888;
  margin-right: 6px;
}
.a4-address {
  margin-top: 8px;
  font-size: 13px;
}
.a4-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.a4-table th,
.a4-table td {
  border: 1px solid #e5e5e5;
  padding: 6px 8px;
  text-align: left;
}
.a4-table th {
  background: #fafafa;
}

@media print {
  .a4-toolbar {
    display: none;
  }
  .a4-page {
    box-shadow: none;
    margin: 0;
  }
}
</style>
