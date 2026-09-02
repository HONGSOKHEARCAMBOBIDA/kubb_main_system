<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { getAdmission, updateAdmission } from '../../services/admission.service'
import { invoicecreate } from '../../services/invoice.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppForm from '@/components/forms/AppForm.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import { studentTermCreate, studentTermUpdate } from '../../services/studentterm.service.js'
import { getSemesterByAcademic } from '../../services/semester.service.js'
import { getAcademics } from '../../services/academic.service.js'
import {
  EnrollmentCreate,
  EnrollmentUpdate,
  EnrollmentDelete,
} from '../../services/enrollment.service.js'
import { getSchoolarshipGroup } from '../../services/schoolarship.service.js'
import { getGenerationByAcademic } from '../../services/generation.service.js'
import { getTermByGeneation } from '../../services/term.service'

const url = '/logo.png'
const notify = useNotification()
let searchTimer = null
const admissions = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

/* ---------------------------------------------------
 * Khmer number-to-words helper (for the receipt)
 * --------------------------------------------------- */
const KH_ONES = ['សូន្យ', 'មួយ', 'ពីរ', 'បី', 'បួន', 'ប្រាំ', 'ប្រាំមួយ', 'ប្រាំពីរ', 'ប្រាំបី', 'ប្រាំបួន']
const KH_TENS = ['', 'ដប់', 'ម្ភៃ', 'សាមសិប', 'សែសិប', 'ហាសិប', 'ហុកសិប', 'ចិតសិប', 'ប៉ែតសិប', 'កៅសិប']

function twoDigitKh(n) {
  if (n === 0) return ''
  if (n < 10) return KH_ONES[n]
  const t = Math.floor(n / 10)
  const o = n % 10
  if (t === 1) {
    return o === 0 ? 'ដប់' : o === 1 ? 'ដប់មួយ' : 'ដប់' + KH_ONES[o]
  }
  let s = KH_TENS[t]
  if (o === 0) return s
  if (o === 1) return s + 'មួយ'
  return s + KH_ONES[o]
}

function integerToKhmerWords(numIn) {
  let n = Math.floor(numIn)
  if (n === 0) return KH_ONES[0]

  const million = Math.floor(n / 1000000); n %= 1000000
  const hunThousand = Math.floor(n / 100000); n %= 100000
  const tenThousand = Math.floor(n / 10000); n %= 10000
  const thousand = Math.floor(n / 1000); n %= 1000
  const hundred = Math.floor(n / 100); n %= 100
  const rest = n

  let parts = []
  if (million > 0) parts.push(integerToKhmerWords(million) + 'លាន')
  if (hunThousand > 0) parts.push(KH_ONES[hunThousand] + 'សែន')
  if (tenThousand > 0) parts.push(KH_ONES[tenThousand] + 'ម៉ឺន')
  if (thousand > 0) parts.push((thousand < 10 ? KH_ONES[thousand] : integerToKhmerWords(thousand)) + 'ពាន់')
  if (hundred > 0) parts.push(KH_ONES[hundred] + 'រយ')
  if (rest > 0) parts.push(twoDigitKh(rest))

  return parts.join('')
}

function amountToKhmerWords(amount) {
  const num = Number(amount) || 0
  const dollars = Math.floor(num)
  const cents = Math.round((num - dollars) * 100)

  let words = integerToKhmerWords(dollars) + 'ដុល្លារ'
  if (cents > 0) {
    words += ' ' + integerToKhmerWords(cents) + 'សេន'
  }
  return words
}
/* --------------------------------------------------- */

const showA4visible = ref(false)
const datashowA4 = ref(null)

// row = admission row, enrollment = enrollment row, feeRow = fee_response row, installmentRow = installment row
async function openshowA4(admission, enrollment, feeRow, installmentRow) {
  datashowA4.value = {
    invoice_code: installmentRow.invoice_code,
    invoice_date: installmentRow.invoice_date,

    student_name: admission.student_name_kh,
    student_name_en: admission.student_name_en,
    student_code: admission.student_code,
    study_level: admission.programme_name,
    generation: admission.generation_name,
    major: admission.major_name,
    course: admission.term_name,

    amount: formatMoney(installmentRow.invoice_grant_total),
    amount_unit: 'ដុល្លារ',
    amount_in_words: amountToKhmerWords(installmentRow.invoice_grant_total),

    payment_method: installmentRow.payment_method,
    payment_type: labelOf(feeIntervalLabels, enrollment.fee_interval),
    payment_number: installmentRow.sequence_no,
  }
  showA4visible.value = true
}

function printA4() {
  window.print()
}

const invoiceDialogVisible = ref(false)
const invoiceSubmitting = ref(false)
const invoiceFormRef = ref(null)

const getStatusLabel = (status) => {
  const labels = {
    PENDING: 'កំពុងសិក្សា',
    FINISH: 'បានបញ្ចប់',
    SUSPENDED: 'ផ្អាកការសិក្សា',
  }
  return labels[status] ?? ''
}

// add new student term
const studyyearOption = [
  { label: 'ឆ្នាំទី1', value: 1 },
  { label: 'ឆ្នាំទី2', value: 2 },
  { label: 'ឆ្នាំទី3', value: 3 },
  { label: 'ឆ្នាំទី4', value: 4 },
]

const AcademicID = ref(null)
const academicOptions = ref([])
const semesterOptions = ref([])
async function fetchAcademicOption() {
  try {
    const res = await getAcademics()
    academicOptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }))
  } catch (e) { }
}
async function onSemesterAcademicChange() {
  semesterOptions.value = await loadSemester(AcademicID.value)
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
const formstudentterm = reactive({
  enrollment_id: null,
  semester_id: null,
  study_year_id: null,
})
const addstudenttermvisible = ref(false)
function openCreateStudentTerm(enrollment) {
  formstudentterm.enrollment_id = enrollment.id
  addstudenttermvisible.value = true
}

function closeCreateStudentTerm() {
  addstudenttermvisible.value = false
}

async function submitstudentterm() {
  try {
    await studentTermCreate(formstudentterm)
    addstudenttermvisible.value = false
    formstudentterm.enrollment_id = null
    fetchAdmissions()
  } catch (e) { }
}
//student term

// update studentterm
const enrollmentStateOptions = [
  { label: 'កំពុងរៀន', value: 'PENDING' },
  { label: 'បានរៀនចប់', value: 'FINISH' },
  { label: 'ព្យួរ / ផ្អាកការសិក្សា', value: 'SUSPENDED' },
]
const AcademicIDUpdateStudentTerm = ref(null)
const updatestudenttermform = reactive({
  semester_id: null,
  study_year_id: null,
  status: '',
})
const studenttermupdateuuid = ref(null)
const updatestudenttermvisible = ref(false)
async function openstudenttermupdate(row) {
  AcademicIDUpdateStudentTerm.value = row.academic_id
  studenttermupdateuuid.value = row.uuid
  updatestudenttermform.semester_id = row.semester_id
  updatestudenttermform.study_year_id = row.study_year_id
  updatestudenttermform.status = row.status
  updatestudenttermvisible.value = true
}
async function closestudenttermupdate(row) {
  studenttermupdateuuid.value = null
  updatestudenttermvisible.value = false
}
async function submitUpdateStudentTerm() {
  try {
    await studentTermUpdate(studenttermupdateuuid.value, {
      semester_id: updatestudenttermform.semester_id,
      study_year_id: updatestudenttermform.study_year_id,
      status: updatestudenttermform.status,
    })
    notify.success('កែប្រែដោយជោគជ័យ')
    updatestudenttermvisible.value = false
    fetchAdmissions()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to update admission')
  }
}
// update studentterm

const filters = reactive({
  student_id: '',
  student_name: '',
  academic_id: null,
  generation_id: null,
  term_id: null,
})

const updateenrollmentvisible = ref(false)
const enrollmentuuid = ref(null)
const formupdateenrollment = reactive({
  description: '',
})

async function openenrollment(row) {
  enrollmentuuid.value = row.uuid
  formupdateenrollment.description = row.description
  updateenrollmentvisible.value = true
}

async function closeenrollment(row) {
  updateenrollmentvisible.value = false
}

async function submitUpdateEnrollment() {
  try {
    await EnrollmentUpdate(enrollmentuuid.value, {
      description: formupdateenrollment.description,
    })
    notify.success('កែប្រែដោយជោគជ័យ')
    updateenrollmentvisible.value = false
    fetchAdmissions()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to update admission')
  }
}

const form = reactive({
  installment_uuid: '',
  fee_id: null,
  invoice_date: '',
  due_date: '',
  total: 0,
  discount: 0,
  tax: 0,
  grant_total: 0,
  message_on_invoice: '',
  reference: '',
  method: '',
})

// Enrollment
const schoolarshipOptions = ref([])
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
const endrollmentFeeintervalOption = [
  { label: 'បង់1ខែម្ដង', value: 'monthly_fee' },
  { label: 'បង់3ខែម្ដង', value: 'quarterly_fee' },
  { label: 'បង់1ឆមាសម្ដង', value: 'semesterly_fee' },
  { label: 'បង់1ឆ្នាំម្ដង', value: 'yearly_fee' },
]
const formenrollment = reactive({
  admision_id: null,
  scholarship_id: null,
  fee_interval: '',
  student_term: {
    semester_id: null,
    study_year_id: null,
  },
})

// Filter
const generationFilterOptions = ref([])
const termFilterOptions = ref([])
const academicfilter = ref(null)
const generationfilter = ref(null)
const termfilter = ref(null)
// Filter

// formupdateadmission
const updateAdmissionVisible = ref(false)
const admissionUpdateID = ref(null)
const selectacademicforupdateadmission = ref(null)
const selectgenerationforupdateadmission = ref(null)
const selecttermforupdateadmission = ref(null)
const generationOptions = ref([])
const termOptions = ref([])
const formupdateadmission = reactive({
  term_id: null,
  state: '',
})
async function openUpdateAdmission(row) {
  admissionUpdateID.value = row.uuid
  formupdateadmission.state = row.state
  formupdateadmission.term_id = row.term_id ?? null

  selectacademicforupdateadmission.value = row.academic_id ?? null
  generationOptions.value = selectacademicforupdateadmission.value
    ? await loadGeneration(selectacademicforupdateadmission.value)
    : []

  selectgenerationforupdateadmission.value = row.generation_id ?? null
  termOptions.value = selectgenerationforupdateadmission.value
    ? await loadTerm(selectgenerationforupdateadmission.value)
    : []

  updateAdmissionVisible.value = true
}

function closeUpdateAdmission() {
  updateAdmissionVisible.value = false
  admissionUpdateID.value = null
}

async function submitUpdateAdmission() {
  try {
    await updateAdmission(admissionUpdateID.value, {
      term_id: formupdateadmission.term_id,
      state: formupdateadmission.state,
    })
    notify.success('កែប្រែពាក្យសុំដោយជោគជ័យ')
    updateAdmissionVisible.value = false
    fetchAdmissions()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to update admission')
  }
}

const admissionStateOptions = [
  { label: 'កំពុងរង់ចាំ', value: 'created' },
  { label: 'ត្រូវបានទទួលយក', value: 'approved' },
  { label: 'បដិសេធ', value: 'rejected' },
]

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

async function onAcademicFilterChange() {
  filters.generation_id = null
  filters.term_id = null
  generationFilterOptions.value = await loadGeneration(filters.academic_id)
  termFilterOptions.value = []
}

async function onGenerationFilterChange() {
  filters.term_id = null
  termFilterOptions.value = await loadTerm(filters.generation_id)
}

async function onAcademicChange() {
  selectgenerationforupdateadmission.value = null
  generationOptions.value = await loadGeneration(selectacademicforupdateadmission.value)
}

async function onGenerationChange() {
  selecttermforupdateadmission.value = null
  termOptions.value = await loadTerm(selectgenerationforupdateadmission.value)
}

// formupdateadmission

// DeleteEnrollment
const deleteenrollmentvisible = ref(false)
const enrollmentuuidfordelete = ref(null)
async function openDeleteEnrollment(row) {
  enrollmentuuidfordelete.value = row.uuid
  deleteenrollmentvisible.value = true
}

function closedeleteenrollment() {
  deleteenrollmentvisible.value = false
}

async function submitDeleteEnrollment() {
  try {
    await EnrollmentDelete(enrollmentuuidfordelete.value)
    notify.success('លុបបានដោយជោគជ័យ')
    deleteenrollmentvisible.value = false
    fetchAdmissions()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to update admission')
  }
}
// DeleteEnrollment

const addenrollmentvisible = ref(false)
function openCreateEnrollment(admission) {
  formenrollment.admision_id = admission.id
  AcademicID.value = admission.academic_id

  addenrollmentvisible.value = true
}
function closeCreateEnrollment() {
  addstudenttermvisible.value = false
}

async function submitEnrollment() {
  try {
    await EnrollmentCreate(formenrollment)
    addenrollmentvisible.value = false
    formenrollment.admision_id = null
    fetchAdmissions()
  } catch (e) { }
}

const invoiceRules = {
  fee_id: [{ required: true, message: 'សូមជ្រើសរើសថ្លៃសិក្សា', trigger: 'change' }],
  invoice_date: [{ required: true, message: 'សូមបញ្ចូលថ្ងៃចេញ', trigger: 'change' }],
  due_date: [{ required: true, message: 'សូមបញ្ចូលថ្ងៃកំណត់', trigger: 'change' }],
  total: [{ required: true, message: 'សូមបញ្ចូលចំនួនទឹកប្រាក់', trigger: 'blur' }],
  grant_total: [{ required: true, message: 'ចំនួនត្រូវបង់មិនអាចទទេ', trigger: 'blur' }],
  method: [{ required: true, message: 'សូមជ្រើសរើសមធ្យោបាយបង់ប្រាក់', trigger: 'change' }],
}

const methodOptions = [
  { label: 'សាច់ប្រាក់', value: 'cash' },
  { label: 'ACELEDA', value: 'aceleda' },
  { label: 'ABA', value: 'aba' },
  { label: 'WING', value: 'wing' },
]

function todayStr() {
  return new Date().toISOString().slice(0, 10)
}

function resetForm() {
  form.installment_uuid = ''
  form.fee_id = null
  form.invoice_date = todayStr()
  form.due_date = todayStr()
  form.total = 0
  form.discount = 0
  form.tax = 0
  form.grant_total = 0
  form.message_on_invoice = ''
  form.reference = ''
  form.method = ''
}

watch(
  () => [form.total, form.discount, form.tax],
  ([t, d, x]) => {
    const g = Number(t || 0) - Number(d || 0) + Number(x || 0)
    form.grant_total = Math.max(0, Number(g.toFixed(2)))
  }
)

function openInvoiceDialog(feeRow, installmentRow) {
  resetForm()
  form.installment_uuid = installmentRow.uuid
  form.fee_id = feeRow?.id ?? feeRow?.ID ?? null
  form.total = Number(installmentRow?.amount ?? feeRow?.total ?? 0)
  form.due_date = installmentRow.due_date
  form.discount = 0
  form.tax = 0
  form.reference = installmentRow ? `លេខយោង-${installmentRow.id}` : ''
  invoiceDialogVisible.value = true
}

function closeInvoiceDialog() {
  invoiceDialogVisible.value = false
  invoiceFormRef.value?.clearValidate?.()
}

async function submitInvoice() {
  if (!invoiceFormRef.value) return
  try {
    await invoiceFormRef.value.validate()
  } catch {
    return
  }

  invoiceSubmitting.value = true
  try {
    const payload = {
      installment_uuid: form.installment_uuid,
      fee_id: Number(form.fee_id),
      invoice_date: form.invoice_date,
      due_date: form.due_date,
      total: Number(form.total),
      discount: Number(form.discount),
      tax: Number(form.tax),
      grant_total: Number(form.grant_total),
      message_on_invoice: form.message_on_invoice,
      reference: form.reference,
      method: form.method,
    }

    await invoicecreate(payload)

    notify.success('បង់ប្រាក់ និងបង្កើតវិក័យបត្រដោយជោគជ័យ')
    invoiceDialogVisible.value = false
    await fetchAdmissions()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to create invoice')
  } finally {
    invoiceSubmitting.value = false
  }
}

// --- Label maps ---
const stateLabels = {
  pending: 'កំពុងរង់ចាំ',
  approved: 'អនុម័ត',
  rejected: 'បដិសេធ',
}
const discountTypeLabels = {
  percentage: 'ភាគរយ',
  amount: 'ចំនួនទឹកប្រាក់',
}
const feeIntervalLabels = {
  monthly_fee: 'ប្រចាំខែ',
  quarterly_fee: 'ប្រចាំត្រីមាស',
  semesterly_fee: 'ប្រចាំឆមាស',
  yearly_fee: 'ប្រចាំឆ្នាំ',
}
const installmentStatusLabels = {
  pending: 'មិនទាន់បង់',
  paid: 'បានបង់',
  overdue: 'ហួសកំណត់',
}

function labelOf(map, value) {
  return map[value] ?? value ?? '-'
}

function formatMoney(v) {
  if (v === null || v === undefined || v === '') return '-'
  return Number(v).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const columns = [
  { slot: 'student_name', label: 'ឈ្មោះ', minwidth: 160 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', slot: 'academic', width: 130 },
  { prop: 'generation_name', label: 'ជំនាន់', slot: 'generation', width: 130 },
  { prop: 'term_name', label: 'វគ្គ', width: 80 },
  { slot: 'major_name', label: 'ជំនាញ', minwidth: 150 },
  { prop: 'programme_name', label: 'កម្រិត', width: 110 },
  { slot: 'group', label: 'ក្រុមបញ្ចុះតម្លៃ', minwidth: 150 },
  { prop: 'date', label: 'ថ្ងៃដាក់ពាក្យ', width: 100 },
  { label: 'ស្ថានភាពពាក្យ', slot: 'state', width: 120 },
  { label: 'ចំនួនដាក់ពាក្យ', slot: 'enrollCount', width: 130 },
]

const columnenrollments = [
  { slot: 'year_id', label: 'ឆ្នាំទី', minwidth: 60 },
  { slot: 'schoolarship', label: 'អាហារូបករណ៍ទទួលបាន', minwidth: 160 },
  { slot: 'fee_interval', label: 'សុំបង់ប្រាក់ជា', minwidth: 130 },
  { slot: 'description', label: 'ផ្សេងៗ' },
]

const columnstudentterms = [
  { slot: 'semester', label: 'ឆមាស', minwidth: 160 },
  { prop: 'study_year_id', label: 'ឆ្នាំទី', minwidth: 160 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 130 },
  { prop: 'status', label: 'ស្ថានភាព', slot: 'stStatus', minwidth: 120 },
]

const columngparecord = [
  { prop: 'total_credit', label: 'Total Credit', minwidth: 160 },
  { prop: 'semester_gpa', label: 'Semester GPA', minwidth: 130 },
  { prop: 'cumulative_gpa', label: 'Cumulative GPA', minwidth: 120 },
]

const columnfees = [
  { prop: 'date', label: 'ថ្ងៃបង្កេីត', minwidth: 110 },
  { label: 'តម្លៃដេីម', slot: 'feeAmount', minwidth: 100 },
  { label: 'បញ្ចុះតម្លៃ', slot: 'feeDiscount', minwidth: 100 },
  { label: 'តម្លៃត្រូវបង់', slot: 'feeTotal', minwidth: 100 },
  { label: 'ស្ថានភាព', slot: 'feeActive', minwidth: 90 },
]

const columninvoices = [
  { prop: 'code', label: 'លេខកូដ', minwidth: 110 },
  { prop: 'invoice_date', label: 'ថ្ងៃចេញ', minwidth: 110 },
  { prop: 'due_date', label: 'ថ្ងៃកំណត់', minwidth: 110 },
  { label: 'សរុប', slot: 'invGrantTotal', minwidth: 100 },
  { label: 'ពន្ធ', slot: 'invTax', minwidth: 90 },
]

const columninstallments = [
  { prop: 'sequence_no', label: 'លេីកទី', width: 70 },
  { prop: 'due_date', label: 'ថ្ងៃត្រូវបង់', width: 110 },
  { label: 'ចំនួនត្រូវបង់', slot: 'instAmount', width: 100 },
  { label: 'ស្ថានភាព', slot: 'instStatus', width: 110 },
  { prop: 'invoice_code', label: 'លេខកូដវិក័យបត្រ', minwidth: 110 },
  { prop: 'invoice_date', label: 'ថ្ងៃចេញវិក័យបត្រ', width: 120 },
  { slot: 'invoice_grant_total', label: 'ចំនួនទទួលបានសរុប', width: 140 },
  { prop: 'payment_code', label: 'លេខកូដទូទាត់', minwidth: 110 },
  { prop: 'payment_reference', label: 'លេខយោង', minwidth: 110 },
  { slot: 'payment_method', label: 'វិធីសាស្រ្តទូទាត់', width: 120 },
  { prop: 'message_on_invoice', label: 'ផ្សេងៗ', minwidth: 120 },
]

const columnpayments = [
  { prop: 'code', label: 'លេខកូដ', minwidth: 110 },
  { prop: 'date', label: 'ថ្ងៃ', minwidth: 110 },
  { label: 'ចំនួន', slot: 'payAmount', minwidth: 100 },
  { prop: 'method', label: 'មធ្យោបាយ', minwidth: 110 },
  { prop: 'reference', label: 'លេខយោង', minwidth: 120 },
]

async function fetchAdmissions() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.student_id) params.student_id = filters.student_id
    if (filters.student_name) params.student_name = filters.student_name
    if (filters.academic_id) params.academic_id = filters.academic_id
    if (filters.generation_id) params.generation_id = filters.generation_id
    if (filters.term_id) params.term_id = filters.term_id
    const res = await getAdmission(params)
    admissions.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load admissions')
  } finally {
    loading.value = false
  }
}

watch(
  () => AcademicIDUpdateStudentTerm.value,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(async () => {
      semesterOptions.value = []
      semesterOptions.value = await loadSemester(AcademicIDUpdateStudentTerm.value)
    }, 400)
  }
)

watch(
  () => filters.generation_id,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchAdmissions()
    }, 400)
  }
)

watch(
  () => filters.term_id,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchAdmissions()
    }, 400)
  }
)

watch(
  () => filters.academic_id,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchAdmissions()
    }, 400)
  }
)

watch(
  () => filters.student_id,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchAdmissions()
    }, 400)
  }
)

watch(
  () => filters.student_name,
  () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      page.value = 1
      fetchAdmissions()
    }, 400)
  }
)

onMounted(() => {
  fetchAdmissions()
  fetchAcademicOption()
  fetchSchoolarship()
})
</script>

<template>
  <div class="admission-page">
    <AppFilterBar :fields="[
      { slot: 'studentId', span: 3 },
      { slot: 'studentName', span: 3 },
      { slot: 'academic', span: 3 },
      { slot: 'generation', span: 3 },
      { slot: 'term', span: 3 },
    ]">
      <template #studentId>
        <AppInput v-model="filters.student_id" placeholder="លេខសម្គាល់សិស្ស" clearable />
      </template>
      <template #studentName>
        <AppInput v-model="filters.student_name" placeholder="ស្វែងរកតាមឈ្មោះសិស្ស" clearable />
      </template>
      <template #academic>
        <AppSelect v-model="filters.academic_id" :options="academicOptions" @change="onAcademicFilterChange" clearable
          placeholder="ឆ្នាំសិក្សា" />
      </template>
      <template #generation>
        <AppSelect v-model="filters.generation_id" :options="generationFilterOptions" @change="onGenerationFilterChange"
          clearable placeholder="ជំនាន់" />
      </template>
      <template #term>
        <AppSelect v-model="filters.term_id" :options="termFilterOptions" placeholder="វគ្គ" clearable @change="" />
      </template>
    </AppFilterBar>

    <TableCustom expandable :data="admissions" :columns="columns" :loading="loading" :total="total"
      v-model:current-page="page" v-model:page-size="pageSize" @page-change="fetchAdmissions">
      <template #student_name="{ row }">
        <el-text>
          <div>
            <el-text size="small">{{ row.student_name_kh }} |
              <el-text size="small">{{ row.student_gender }}</el-text></el-text>
          </div>
          <div>
            <el-text size="small" tag="b" type="primary">{{ row.student_name_en }} |
              <el-text size="small">{{ row.student_code }}</el-text></el-text>
          </div>
        </el-text>
      </template>

      <template #generation="{ row }">
        <el-text>{{ row.generation_name || '-' }}</el-text>
      </template>

      <template #academic="{ row }">
        <el-text>{{ row.academic_name || '-' }}</el-text>
      </template>

      <template #state="{ row }">
        <el-text>{{ labelOf(stateLabels, row.state) }}</el-text>
      </template>

      <template #enrollCount="{ row }">
        <el-text> {{ (row.enrollment || []).length }} ដង </el-text>
      </template>

      <template #major_name="{ row }">
        <el-text>
          <div>
            <el-text>{{ row.major_name }} | {{ row.major_code }}</el-text>
          </div>
          <div>
            <el-text size="small" type="primary">{{ row.yearly_fee }}$ /ឆ្នាំ | {{ row.semesterly_fee }}$ /ឆមាស |
              {{ row.quarterly_fee }}$ /ត្រីមាស | {{ row.monthly_fee }}$ /ខែ</el-text>
          </div>
        </el-text>
      </template>

      <template #group="{ row }">
        <el-text style="color: black">
          {{
            row.discount_type === 'percentage'
              ? `${row.group_name} - បញ្ចុះតម្លៃ ${row.discount_percentage} %`
              : `${row.group_name} - បញ្ចុះតម្លៃ ${row.discount_amount} $`
          }}
        </el-text>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" plain type="primary" @click="openUpdateAdmission(row)" />
        </el-tooltip>
      </template>

      <!-- Level 1: admission -> enrollments -->
      <template #expand="{ row: admissionRow }">
        <el-divider content-position="left">ព័ត៌មានដាក់ពាក្យ
          <AppButton plain @click="openCreateEnrollment(admissionRow)" type="primary">
            ដាក់ពាក្យថ្មី
          </AppButton>
        </el-divider>
        <TableCustom expandable :data="admissionRow.enrollment" :columns="columnenrollments" :show-pagination="false">
          <template #year_id="{ row }">
            <el-text> ឆ្នាំទី {{ row.year_id }} </el-text>
          </template>
          <template #schoolarship="{ row }">
            <el-text v-if="row.schoolarship_id">
              {{ row.schoolarship_name }} |
              <el-text type="primary">
                {{
                  row.schoolarship_discount_type === 'percentage'
                    ? `ទទួលការបញ្ចុះតម្លៃ${row.schoolarship_discount_percentage}%`
                    : `ទទួលការបញ្ចុះតម្លៃ${row.schoolarship_discount_amount}$`
                }}
              </el-text>
            </el-text>
            <el-text v-else type="info">គ្មាន</el-text>
          </template>

          <template #fee_interval="{ row }">
            <el-text>{{ labelOf(feeIntervalLabels, row.fee_interval) }}</el-text>
          </template>

          <template #description="{ row }">
            <el-text>{{ row.description || '-' }}</el-text>
          </template>

          <template #actions="{ row }">
            <el-tooltip content="កែប្រែ" placement="top">
              <AppButton size="small" circle plain icon="Edit" type="primary" @click="openenrollment(row)">
              </AppButton>
            </el-tooltip>
            <el-tooltip content="លុប" placement="top">
              <AppButton size="small" circle plain icon="Delete" type="danger" @click="openDeleteEnrollment(row)">
              </AppButton>
            </el-tooltip>
          </template>

          <!-- Level 2: enrollment -> student_term / fee_response -->
          <template #expand="{ row: enrollmentRow }">
            <el-divider content-position="left">
              កំពុងសិក្សា
              <AppButton plain v-if="enrollmentRow.student_term?.length < 2" @click="openCreateStudentTerm(enrollmentRow)">
                ថែមឆមាសថ្មី
              </AppButton>
            </el-divider>
            <TableCustom expandable :data="enrollmentRow.student_term" :columns="columnstudentterms" :show-pagination="false">
              <template #semester="{ row }">
                <el-text class="ml-2">{{ row.semester_name }}</el-text>
              </template>
              <template #stStatus="{ row }">
                <el-text>
                  {{ getStatusLabel(row.status) }}
                </el-text>
              </template>
              <template #actions="{ row }">
                <el-tooltip content="កែប្រែ" placement="top">
                  <AppButton size="small" circle plain icon="Edit" type="primary" @click="openstudenttermupdate(row)">
                  </AppButton>
                </el-tooltip>
              </template>
              <!-- Level 3: student_term -> gpa_record -->
              <template #expand="{ row: studentTermRow }">
                <el-divider content-position="left"> GPA Record </el-divider>
                <TableCustom expandable :data="studentTermRow.gpa_record" :columns="columngparecord" :show-pagination="false">
                </TableCustom>
              </template>
            </TableCustom>

            <el-divider content-position="left">ថ្លៃសិក្សាប្រចាំឆ្នាំ</el-divider>
            <TableCustom expandable :data="enrollmentRow.fee_response" :columns="columnfees" :show-pagination="false">
              <template #feeAmount="{ row }">
                <el-text style="color: black"> {{ formatMoney(row.amount) }}$ </el-text>
              </template>
              <template #feeDiscount="{ row }">
                <el-text style="color: crimson">{{ formatMoney(row.discount) }}$</el-text>
              </template>
              <template #feeTotal="{ row }">
                <el-text>{{ formatMoney(row.total) }}$</el-text>
              </template>
              <template #feeActive="{ row }">
                <el-tag :type="row.active ? 'danger' : 'info'" size="small">
                  {{ row.active ? 'បង់មិនទាន់គ្រប់' : 'បង់គ្រប់' }}
                </el-tag>
              </template>

              <!-- Level 4: fee_response -> installment -->
              <template #expand="{ row: feeRow }">
                <div>
                  <el-divider content-position="left"> ការបង់ប្រាក់ </el-divider>
                  <TableCustom expandable :data="feeRow.installment || []" :columns="columninstallments"
                    :show-pagination="false" actions-width="180">
                    <template #instAmount="{ row }">
                      <el-text>{{ formatMoney(row.amount) }}$ </el-text>
                    </template>
                    <template #instStatus="{ row }">
                      <el-text :type="row.status === 'pending' ? 'danger' : 'success'">
                        {{ labelOf(installmentStatusLabels, row.status) }}
                      </el-text>
                    </template>
                    <template #invoice_grant_total="{ row }">
                      <el-text>{{ formatMoney(row.invoice_grant_total) }}$</el-text>
                    </template>
                    <template #payment_method="{ row }">
                      <el-text style="color: black">{{ row.payment_method }}</el-text>
                    </template>
                    <template #actions="{ row: installmentRow }">
                      <AppButton :disabled="installmentRow.status === 'paid'" type="success"
                        @click="openInvoiceDialog(feeRow, installmentRow)">
                        បង់ប្រាក់
                      </AppButton>
                      <AppButton :disabled="installmentRow.status === 'pending'" type="primary"
                        @click="openshowA4(admissionRow, enrollmentRow, feeRow, installmentRow)">
                        បោះពុម្ព
                      </AppButton>
                    </template>
                  </TableCustom>
                </div>
              </template>
            </TableCustom>
          </template>
        </TableCustom>
      </template>
    </TableCustom>

    <AppDialog v-if="invoiceDialogVisible" v-model:visible="invoiceDialogVisible" title="បង្កើតវិក័យបត្រ / បង់ប្រាក់"
      :showDefaultFooter="false" width="520px" @close="closeInvoiceDialog">
      <AppForm ref="invoiceFormRef" :model="form" :rules="invoiceRules" :show-actions="true" @submit="submitInvoice"
        submitText="បញ្ជាក់ការបង់ប្រាក់">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="ថ្ងៃចេញ" prop="invoice_date">
              <el-date-picker v-model="form.invoice_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="ថ្ងៃកំណត់" prop="due_date">
              <el-date-picker v-model="form.due_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput v-model.number="form.total" type="number" label="តម្លៃដេីម" placeholder="តម្លៃដេីម">
            </AppInput>
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="form.discount" type="number" label="បញ្ចុះតម្លៃ" placeholder="បញ្ចុះតម្លៃ">
            </AppInput>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput v-model.number="form.tax" type="number" label="ពន្ធ" placeholder="ពន្ធ">
            </AppInput>
          </el-col>
          <el-col :span="12">
            <AppInput v-model.number="form.grant_total" type="number" label="សរុបត្រូវបង់" placeholder="សរុបត្រូវបង់">
            </AppInput>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect v-model="form.method" :options="methodOptions" placeholder="មធ្យោបាយបង់ប្រាក់"
              label="មធ្យោបាយបង់ប្រាក់" clearable />
          </el-col>
          <el-col :span="12">
            <AppInput v-model="form.reference" label="លេខយោង" placeholder="លេខយោង"> </AppInput>
          </el-col>
        </el-row>

        <el-row>
          <el-col>
            <AppInput type="area" v-model="form.message_on_invoice" label="សេចក្តីលម្អិត" placeholder="សេចក្តីលម្អិត">
            </AppInput>
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>

    <AppDialog v-if="addstudenttermvisible" v-model:visible="addstudenttermvisible" title="ថែមឆមាស"
      :showDefaultFooter="false" width="520px" @close="closeCreateStudentTerm">
      <AppForm :model="formstudentterm" :show-actions="true" @submit="submitstudentterm" submitText="រក្សាទុក">
        <AppSelect v-model="AcademicID" :options="academicOptions" placeholder="រេីសឆ្នាំសិក្សា" label="ឆ្នាំសិក្សា"
          clearable @change="onSemesterAcademicChange" />

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect v-model="formstudentterm.study_year_id" :options="studyyearOption" placeholder="ឆ្នាំទី"
              label="ឆ្នាំទី" clearable />
          </el-col>
          <el-col :span="12">
            <AppSelect v-model="formstudentterm.semester_id" :options="semesterOptions" :disabled="!AcademicID"
              placeholder="ឆមាសទី" label="ឆមាសទី" clearable />
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>
    <AppDialog v-if="addenrollmentvisible" v-model:visible="addenrollmentvisible" title="ថែមការចុះឈ្មោះ"
      :showDefaultFooter="false" width="720px" @close="closeCreateEnrollment">
      <AppForm :model="formenrollment" :show-actions="true" @submit="submitEnrollment" submitText="រក្សាទុក">
        <el-row :gutter="20">
          <el-col :span="8">
            <AppSelect v-model="AcademicID" :options="academicOptions" placeholder="រេីសឆ្នាំសិក្សា" label="ឆ្នាំសិក្សា"
              clearable @change="onSemesterAcademicChange" />
          </el-col>
          <el-col :span="8">
            <AppSelect v-model="formenrollment.scholarship_id" :options="schoolarshipOptions" placeholder="អាហារូបករណ៍"
              label="អាហារូបករណ៍" clearable />
          </el-col>
          <el-col :span="8">
            <AppSelect v-model="formenrollment.fee_interval" :options="endrollmentFeeintervalOption"
              placeholder="សុំបង់ប្រាក់ជា" label="សុំបង់ប្រាក់ជា" clearable />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect v-model="formenrollment.student_term.study_year_id" :options="studyyearOption"
              placeholder="ឆ្នាំទី" label="ឆ្នាំទី" clearable />
          </el-col>
          <el-col :span="12">
            <AppSelect v-model="formenrollment.student_term.semester_id" :options="semesterOptions"
              :disabled="!AcademicID" placeholder="ឆមាសទី" label="ឆមាសទី" clearable />
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>
    <AppDialog v-if="updateAdmissionVisible" v-model:visible="updateAdmissionVisible" title="កែប្រែពាក្យសុំចូលរៀន"
      :showDefaultFooter="false" width="600px" @close="closeUpdateAdmission">
      <AppForm :model="formupdateadmission" :show-actions="true" @submit="submitUpdateAdmission" submitText="រក្សាទុក">
        <el-row :gutter="20">
          <el-col :span="8">
            <AppSelect v-model="selectacademicforupdateadmission" :options="academicOptions" placeholder="ឆ្នាំសិក្សា"
              label="ឆ្នាំសិក្សា" clearable @change="onAcademicChange" />
          </el-col>
          <el-col :span="8">
            <AppSelect v-model="selectgenerationforupdateadmission" :options="generationOptions"
              :disabled="!selectacademicforupdateadmission" placeholder="ជំនាន់" label="ជំនាន់" clearable
              @change="onGenerationChange" />
          </el-col>
          <el-col :span="8">
            <AppSelect v-model="formupdateadmission.term_id" :options="termOptions"
              :disabled="!selectgenerationforupdateadmission" placeholder="វគ្គ" label="វគ្គ" clearable />
          </el-col>
        </el-row>
        <el-row>
          <el-col>
            <AppSelect v-model="formupdateadmission.state" :options="admissionStateOptions" placeholder="ស្ថានភាពពាក្យ"
              label="ស្ថានភាពពាក្យ" clearable />
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>

    <AppDialog v-if="updateenrollmentvisible" v-model:visible="updateenrollmentvisible" title="កែប្រែ"
      :showDefaultFooter="false" width="600px" @close="closeenrollment">
      <AppForm :model="formupdateenrollment" :show-actions="true" @submit="submitUpdateEnrollment"
        submitText="រក្សាទុក">
        <AppInput type="area" v-model="formupdateenrollment.description" label="សេចក្តីលម្អិត"
          placeholder="សេចក្តីលម្អិត">
        </AppInput>
      </AppForm>
    </AppDialog>

    <AppDialog v-if="deleteenrollmentvisible" v-model:visible="deleteenrollmentvisible" title="លុប"
      :showDefaultFooter="false" width="600px" @close="closedeleteenrollment">
      <AppForm :show-actions="true" @submit="submitDeleteEnrollment" submitText="លុប">
        <el-alert type="warning">គ្រប់ព័ត៏មាននិងត្រូវបានលុប!</el-alert>
      </AppForm>
    </AppDialog>

    <AppDialog v-if="updatestudenttermvisible" v-model:visible="updatestudenttermvisible" title="កែប្រែឆមាស"
      :showDefaultFooter="false" width="520px" @close="closestudenttermupdate">
      <AppForm :model="updatestudenttermform" :show-actions="true" @submit="submitUpdateStudentTerm"
        submitText="រក្សាទុក">
        <AppSelect v-model="AcademicIDUpdateStudentTerm" :options="academicOptions" placeholder="រេីសឆ្នាំសិក្សា"
          label="ឆ្នាំសិក្សា" clearable />

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect v-model="updatestudenttermform.study_year_id" :options="studyyearOption" placeholder="ឆ្នាំទី"
              label="ឆ្នាំទី" clearable />
          </el-col>
          <el-col :span="12">
            <AppSelect v-model="updatestudenttermform.semester_id" :options="semesterOptions"
              :disabled="!AcademicIDUpdateStudentTerm" placeholder="ឆមាសទី" label="ឆមាសទី" clearable />
          </el-col>
        </el-row>
        <AppSelect v-model="updatestudenttermform.status" :options="enrollmentStateOptions" placeholder="ស្ថានភាព"
          label="ស្ថានភាព" clearable />
      </AppForm>
    </AppDialog>

    <AppDialog v-if="showA4visible" v-model:visible="showA4visible" :showDefaultFooter="false" width="58%">
      <AppForm :show-actions="true" @submit="printA4" submitText="បោះពុម្ព">
        <div class="a4-preview-wrapper">
          <div class="a4-landscape">

            <!-- Header -->
            <div class="print-header">
              <div class="header-logo">
                <el-image :src="url" fit="contain" class="university-logo" />
              </div>

              <div class="header-info">
                <div class="university-name-kh">សាកលវិទ្យាល័យ ខេមរាវិទូ</div>
                <div class="university-name-en">KHEMARAVITU UNIVERSITY</div>
                <div class="university-branch">សាខាខេត្តបាត់ដំបង</div>
                <div class="university-branch">បង្កាន់ដៃទទួលប្រាក់</div>
                <div class="university-name-en">Receipt</div>
              </div>

              <div class="header-right">
                <div class="document-box">
                  <el-row class="pl-5 pt-2">
                    <el-text>No: {{ datashowA4.invoice_code }}</el-text>
                  </el-row>
                  <el-row class="pl-5 pb-2">
                    <el-text>Date: {{ datashowA4.invoice_date }}</el-text>
                  </el-row>
                </div>
              </div>
            </div>

            <!-- Content -->
            <div class="receipt-content">
              <div class="info-box">
                <div class="info-row">
                  <div class="info-cell">ឈ្មោះ {{ datashowA4.student_name }}</div>
                  <div class="info-cell">លេខសម្គាល់សិស្ស {{ datashowA4.student_code }}</div>
                </div>
                <div class="info-row">
                  <div class="info-cell">កម្រិតសិក្សា {{ datashowA4.study_level }}</div>
                  <div class="info-cell">ជំនាន់ទី {{ datashowA4.generation }}</div>
                </div>
                <div class="info-row">
                  <div class="info-cell">ជំនាញ {{ datashowA4.major }}</div>
                  <div class="info-cell">វគ្គ {{ datashowA4.course }}</div>
                </div>
              </div>

              <div class="info-box" style="margin-top: 12px;">
                <div class="info-row">
                  <div class="info-cell">ចំនួនទឹកប្រាក់ {{ datashowA4.amount }}</div>
                  <div class="info-cell">ក្នុងខ្ទង់ {{ datashowA4.amount_unit }}</div>
                </div>
                <div class="info-row">
                  <div class="info-cell">ជាអក្សរ {{ datashowA4.amount_in_words }}</div>
                  <div class="info-cell">វិធីសាស្រ្តទូទាត់ {{ datashowA4.payment_method }}</div>
                </div>
                <div class="info-row">
                  <div class="info-cell">របៀបបង់ប្រាក់ {{ datashowA4.payment_type }}</div>
                  <div class="info-cell">បង់លើកទី {{ datashowA4.payment_number }}</div>
                </div>
              </div>
            </div>

            <div class="signature-section">
              <div class="signature-row">
                <div class="signature-item"><div class="signature-title">អ្នកបង់ប្រាក់</div></div>
                <div class="signature-item"><div class="signature-title">អ្នកទទួល</div></div>
                <div class="signature-item"><div class="signature-title">ហេរញ្ញិក</div></div>
              </div>
            </div>

            <div class="receipt-footer">
              <div class="footer-address">
                អាសយដ្ឋាន : ភូមិកម្មករ សង្កាត់ស្វាយប៉ោ ក្រុងបាត់ដំបង ខេត្តបាត់ដំបង
                ទូរសព្ទ ០១២ ៨២៥ ២៥៦ / ០៥៣ ៩៥២ ៣០៦
              </div>
              <div class="footer-note">
                បញ្ជាក់ : ទឹកប្រាក់ដែលបានបង់រួច មិនអាចដក ឬបង្វិលទៅឲ្យអ្នកផ្សេងបានទេ ។
                សូមរក្សាបង្កាន់ដៃនេះឲ្យបានល្អ ។
              </div>
            </div>

          </div>
        </div>
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Moul&family=Siemreap&display=swap');

.a4-preview-wrapper {
  width: 100%;
  padding: 20px;
  background: #f5f6f8;
  overflow: auto;
  box-sizing: border-box;
}

.a4-landscape {
  width: 297mm;
  height: 210mm;
  margin: 0 auto;
  padding: 12mm 15mm;
  background: #fff;
  box-sizing: border-box;
  font-family: "Siemreap", sans-serif;
  color: #000;
}

.print-header {
  display: grid;
  grid-template-columns: 100px 1fr 190px;
  align-items: center;
  min-height: 90px;
  padding-bottom: 12px;
}

.header-logo {
  display: flex;
  align-items: center;
  justify-content: center;
}

.university-logo {
  width: 85px;
  height: 85px;
}

.header-info {
  text-align: center;
  line-height: 1.45;
}

.university-name-kh {
  font-family: "Moul", serif;
  font-size: 21px;
  line-height: 1.5;
}

.university-name-en {
  font-size: 14px;
  line-height: 1.5;
  letter-spacing: 0.3px;
}

.university-branch {
  font-family: "Moul", serif;
  font-size: 12px;
  line-height: 1.5;
}

.header-right {
  display: flex;
  justify-content: flex-end;
}

.document-box {
  width: 100%;
  border: 1px solid #222;
  border-radius: 6px;
  box-sizing: border-box;
}

.header-right .document-box {
  min-height: 55px;
}

.receipt-content {
  margin-top: 50px;
}

.info-box {
  border: 1px solid #222;
  border-radius: 6px;
  overflow: hidden;
  height: 100%;
}

.info-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
}

.info-cell {
  min-height: 34px;
  padding: 7px 12px;
  box-sizing: border-box;
  border-bottom: 1px solid #ddd;
}

.info-cell:nth-child(odd) {
  border-right: 1px solid #ddd;
}

.info-row:last-child .info-cell {
  border-bottom: none;
}

.signature-section {
  margin-top: 35px;
}

.signature-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  text-align: center;
}

.signature-item {
  min-height: 80px;
}

.signature-title {
  font-family: "Siemreap", sans-serif;
  font-size: 14px;
}

.receipt-footer {
  margin-top: 60px;
  border: 1px solid #222;
  border-radius: 6px;
  overflow: hidden;
  font-family: "Moul", serif;
  font-size: 11px;
}

.footer-address,
.footer-note {
  padding: 8px 12px;
  line-height: 2;
}

@media print {
  @page {
    size: A4 landscape;
    margin: 0;
  }

  body * {
    visibility: hidden;
  }
  .a4-preview-wrapper{
    visibility: visible;
  }

  /* Pull the receipt out of the dialog flow and pin it to the page */
  .a4-preview-wrapper {
  width: 297mm;
  height: 210mm;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    padding: 0;
    margin: 0;
    background: #fff;
  }

}
</style>