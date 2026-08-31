<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { getAdmission,updateAdmission } from '../../services/admission.service'
import { invoicecreate } from '../../services/invoice.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppForm from '@/components/forms/AppForm.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import { studentTermCreate } from '../../services/studentterm.service.js'
import { getSemesterByAcademic } from '../../services/semester.service.js'
import { getAcademics } from '../../services/academic.service.js'
import { EnrollmentCreate } from '../../services/enrollment.service.js'
import { getSchoolarshipGroup } from '../../services/schoolarship.service.js'
import { getGenerationByAcademic } from '../../services/generation.service.js'
import { getTermByGeneation } from '../../services/term.service'
const notify = useNotification()

const admissions = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const invoiceDialogVisible = ref(false)
const invoiceSubmitting = ref(false)
const invoiceFormRef = ref(null)

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
  } catch (e) {}
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
  } catch (e) {}
}
//student term

const filters = reactive({
  student_id: '',
  student_name: '',
})

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
  state: ''
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

async function onAcademicChange() {
  selectgenerationforupdateadmission.value = null
  generationOptions.value = await loadGeneration(selectacademicforupdateadmission.value)
}

async function onGenerationChange() {
  selecttermforupdateadmission.value = null
  termOptions.value = await loadTerm(selectgenerationforupdateadmission.value)
}

// formupdateadmission

const addenrollmentvisible = ref(false)
function openCreateEnrollment(admission) {
  formenrollment.admision_id = admission.id
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
  } catch (e) {}
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

// --- Label maps: replace with your real enum values ---
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
  { label: 'សកម្ម', slot: 'stActive', minwidth: 90 },
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
  { prop: 'due_date', label: 'ថ្ងៃត្រូវបង់', minwidth: 110 },
  { label: 'ចំនួនត្រូវបង់', slot: 'instAmount', minwidth: 100 },
  { label: 'ស្ថានភាព', slot: 'instStatus', minwidth: 110 },
  { prop: 'invoice_code', label: 'លេខកូដវិក័យបត្រ', minwidth: 110 },
  { prop: 'invoice_date', label: 'ថ្ងៃចេញវិក័យបត្រ', minwidth: 110 },
  { slot: 'invoice_grant_total', label: 'ចំនួនទទួលបានសរុប', minwidth: 110 },
  { prop: 'payment_code', label: 'លេខកូដទូទាត់', minwidth: 110 },
  { prop: 'payment_reference', label: 'លេខយោង', minwidth: 110 },
  { slot: 'payment_method', label: 'វិធីសាស្រ្តទូទាត់', minwidth: 110 },
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

    const res = await getAdmission(params)
    admissions.value = res.data.data || []
    total.value = res.data.pagination.totalCount || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load admissions')
  } finally {
    loading.value = false
  }
}

function onFilterChange() {
  page.value = 1
  fetchAdmissions()
}

onMounted(() => {
  fetchAdmissions()
  fetchAcademicOption()
  fetchSchoolarship()
})
</script>

<template>
  <div class="admission-page">
    <AppFilterBar
      :fields="[
        { slot: 'studentId', span: 5 },
        { slot: 'studentName', span: 5 },
      ]"
      :action-span="3"
    >
      <template #studentId>
        <AppInput
          v-model="filters.student_id"
          placeholder="លេខសម្គាល់សិស្ស"
          clearable
          @change="onFilterChange"
        />
      </template>
      <template #studentName>
        <AppInput
          v-model="filters.student_name"
          placeholder="ស្វែងរកតាមឈ្មោះសិស្ស"
          clearable
          @change="onFilterChange"
        />
      </template>
    </AppFilterBar>

    <TableCustom
      expandable
      :data="admissions"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchAdmissions"
    >
      <template #student_name="{ row }">
        <el-text>
          <div>
            <el-text size="small"
              >{{ row.student_name_kh }} |
              <el-text size="small">{{ row.student_gender }}</el-text></el-text
            >
          </div>
          <div>
            <el-text type="primary">{{ row.student_name_en }}</el-text>
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
            <el-text size="small" type="primary"
              >{{ row.yearly_fee }}$ /ឆ្នាំ | {{ row.semesterly_fee }}$ /ឆមាស |
              {{ row.quarterly_fee }}$ /ត្រីមាស | {{ row.monthly_fee }}$ /ខែ</el-text
            >
          </div>
        </el-text>
      </template>

      >
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
    <AppButton
      icon="Edit"
      circle
      size="small"
      plain
      type="warning"
      @click="openUpdateAdmission(row)"
    />
  </el-tooltip>
</template>

      <!-- Level 1: enrollments -->
      <template #expand="{ row }">
        <el-divider content-position="left"
          >ព័ត៌មានដាក់ពាក្យ
          <AppButton
            v-if="row.enrollment?.length < 5"
            @click="openCreateEnrollment(row)"
            type="primary"
          >
            ដាក់ពាក្យថ្មី
          </AppButton>
        </el-divider>
        <TableCustom
          expandable
          :data="row.enrollment"
          :columns="columnenrollments"
          :show-pagination="false"
        >
          <template #year_id="{ row }">
            <el-text tag="b"> ឆ្នាំទី {{ row.year_id }} </el-text>
          </template>
          <template #schoolarship="{ row }">
            <el-text v-if="row.schoolarship_id" tag="b" style="color: crimson">
              {{ row.schoolarship_name }} —
              {{
                row.schoolarship_discount_type === 'percentage'
                  ? `ទទួលការបញ្ចុះតម្លៃ${row.schoolarship_discount_percentage}%`
                  : `ទទួលការបញ្ចុះតម្លៃ${row.schoolarship_discount_amount}$`
              }}
            </el-text>
            <el-text v-else type="info">គ្មាន</el-text>
          </template>

          <template #fee_interval="{ row }">
            <el-text tag="b" type="success">{{
              labelOf(feeIntervalLabels, row.fee_interval)
            }}</el-text>
          </template>

          <template #description="{ row }">
            <el-text>{{ row.description || '-' }}</el-text>
          </template>

          <template #expand="{ row }">
            <el-divider content-position="left">
              កំពុងសិក្សា
              <AppButton v-if="row.student_term?.length < 2" @click="openCreateStudentTerm(row)">
                ថែមឆមាសថ្មី
              </AppButton>
            </el-divider>
            <TableCustom
              expandable
              :data="row.student_term"
              :columns="columnstudentterms"
              :show-pagination="false"
            >
              <template #semester="{ row }">
                <el-text tag="b" class="ml-2">{{ row.semester_name }}</el-text>
              </template>
              <template #stStatus="{ row }">
                <el-text tag="b" :type="row.status == 'PENDING' ? 'success' : 'info'">
                  {{ row.status }}
                </el-text>
              </template>
              <template #stActive="{ row }">
                <el-tag :type="row.active ? 'success' : 'danger'" size="small">
                  {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
                </el-tag>
              </template>
              <template #expand="{ row }">
                <el-divider content-position="left"> GPA Record </el-divider>
                <TableCustom
                  expandable
                  :data="row.gpa_record"
                  :columns="columngparecord"
                  :show-pagination="false"
                >
                </TableCustom>
              </template>
            </TableCustom>

            <el-divider content-position="left">ថ្លៃសិក្សាប្រចាំឆ្នាំ</el-divider>
            <TableCustom
              expandable
              :data="row.fee_response"
              :columns="columnfees"
              :show-pagination="false"
            >
              <template #feeAmount="{ row }">
                <el-text tag="b" style="color: black"> {{ formatMoney(row.amount) }}$ </el-text>
              </template>
              <template #feeDiscount="{ row }">
                <el-text tag="b" style="color: crimson">{{ formatMoney(row.discount) }}$</el-text>
              </template>
              <template #feeTotal="{ row }">
                <el-text tag="b" type="primary">{{ formatMoney(row.total) }}$</el-text>
              </template>
              <template #feeActive="{ row }">
                <el-tag :type="row.active ? 'success' : 'danger'" size="small">
                  {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
                </el-tag>
              </template>

              <template #expand="{ row: feeRow }">
                <div>
                  <el-divider content-position="left"> ការបង់រំលស់ </el-divider>
                  <TableCustom
                    expandable
                    :data="feeRow.installment || []"
                    :columns="columninstallments"
                    :show-pagination="false"
                    actions-width="180"
                  >
                    <template #instAmount="{ row }">
                      <el-text tag="mark">{{ formatMoney(row.amount) }}$ </el-text>
                    </template>
                    <template #instStatus="{ row }">
                      <el-text tag="b" :type="row.status === 'pending' ? 'danger' : 'success'">
                        {{ labelOf(installmentStatusLabels, row.status) }}
                      </el-text>
                    </template>
                    <template #invoice_grant_total="{ row }">
                      <el-text tag="b">{{ row.invoice_grant_total }}$</el-text>
                    </template>
                    <template #payment_method="{ row }">
                      <el-text tag="b" style="color: black">{{ row.payment_method }}</el-text>
                    </template>
                    <template #actions="{ row: installmentRow }">
                      <AppButton
                        :disabled="installmentRow.status === 'paid'"
                        type="success"
                        @click="openInvoiceDialog(feeRow, installmentRow)"
                      >
                        បង់ប្រាក់
                      </AppButton>
                      <AppButton
                        :disabled="installmentRow.status === 'pending'"
                        type="primary"
                        @click=""
                      >
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
    <AppDialog
      v-if="invoiceDialogVisible"
      v-model:visible="invoiceDialogVisible"
      title="បង្កើតវិក័យបត្រ / បង់ប្រាក់"
      :showDefaultFooter="false"
      width="520px"
      @close="closeInvoiceDialog"
    >
      <AppForm
        ref="invoiceFormRef"
        :model="form"
        :rules="invoiceRules"
        :show-actions="true"
        @submit="submitInvoice"
        submitText="បញ្ជាក់ការបង់ប្រាក់"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="ថ្ងៃចេញ" prop="invoice_date">
              <el-date-picker
                v-model="form.invoice_date"
                type="date"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="ថ្ងៃកំណត់" prop="due_date">
              <el-date-picker
                v-model="form.due_date"
                type="date"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model.number="form.total"
              type="number"
              label="តម្លៃដេីម"
              placeholder="តម្លៃដេីម"
            >
            </AppInput>
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.discount"
              type="number"
              label="បញ្ចុះតម្លៃ"
              placeholder="បញ្ចុះតម្លៃ"
            >
            </AppInput>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput v-model.number="form.tax" type="number" label="ពន្ធ" placeholder="ពន្ធ">
            </AppInput>
          </el-col>
          <el-col :span="12">
            <AppInput
              v-model.number="form.grant_total"
              type="number"
              label="សរុបត្រូវបង់"
              placeholder="សរុបត្រូវបង់"
            >
            </AppInput>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="form.method"
              :options="methodOptions"
              placeholder="មធ្យោបាយបង់ប្រាក់"
              label="មធ្យោបាយបង់ប្រាក់"
              clearable
            />
          </el-col>
          <el-col :span="12">
            <AppInput v-model="form.reference" label="លេខយោង" placeholder="លេខយោង"> </AppInput>
          </el-col>
        </el-row>

        <el-row>
          <el-col>
            <AppInput
              type="area"
              v-model="form.message_on_invoice"
              label="សេចក្តីលម្អិត"
              placeholder="សេចក្តីលម្អិត"
            >
            </AppInput>
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>

    <AppDialog
      v-if="addstudenttermvisible"
      v-model:visible="addstudenttermvisible"
      title="ថែមឆមាស"
      :showDefaultFooter="false"
      width="520px"
      @close="closeCreateStudentTerm"
    >
      <AppForm
        :model="formstudentterm"
        :show-actions="true"
        @submit="submitstudentterm"
        submitText="រក្សាទុក"
      >
        <AppSelect
          v-model="AcademicID"
          :options="academicOptions"
          placeholder="រេីសឆ្នាំសិក្សា"
          label="ឆ្នាំសិក្សា"
          clearable
          @change="onSemesterAcademicChange"
        />

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="formstudentterm.study_year_id"
              :options="studyyearOption"
              placeholder="ឆ្នាំទី"
              label="ឆ្នាំទី"
              clearable
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="formstudentterm.semester_id"
              :options="semesterOptions"
              :disabled="!AcademicID"
              placeholder="ឆមាសទី"
              label="ឆមាសទី"
              clearable
            />
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>
    <AppDialog
      v-if="addenrollmentvisible"
      v-model:visible="addenrollmentvisible"
      title="ថែមការចុះឈ្មោះ"
      :showDefaultFooter="false"
      width="720px"
      @close="closeCreateEnrollment"
    >
      <AppForm
        :model="formenrollment"
        :show-actions="true"
        @submit="submitEnrollment"
        submitText="រក្សាទុក"
      >
        <el-row :gutter="20">
          <el-col :span="8">
            <AppSelect
              v-model="AcademicID"
              :options="academicOptions"
              placeholder="រេីសឆ្នាំសិក្សា"
              label="ឆ្នាំសិក្សា"
              clearable
              @change="onSemesterAcademicChange"
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="formenrollment.scholarship_id"
              :options="schoolarshipOptions"
              placeholder="អាហារូបករណ៍"
              label="អាហារូបករណ៍"
              clearable
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="formenrollment.fee_interval"
              :options="endrollmentFeeintervalOption"
              placeholder="សុំបង់ប្រាក់ជា"
              label="សុំបង់ប្រាក់ជា"
              clearable
            />
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <AppSelect
              v-model="formenrollment.student_term.study_year_id"
              :options="studyyearOption"
              placeholder="ឆ្នាំទី"
              label="ឆ្នាំទី"
              clearable
            />
          </el-col>
          <el-col :span="12">
            <AppSelect
              v-model="formenrollment.student_term.semester_id"
              :options="semesterOptions"
              :disabled="!AcademicID"
              placeholder="ឆមាសទី"
              label="ឆមាសទី"
              clearable
            />
          </el-col>
        </el-row>
      </AppForm>
    </AppDialog>
<AppDialog
  v-if="updateAdmissionVisible"
  v-model:visible="updateAdmissionVisible"
  title="កែប្រែពាក្យសុំចូលរៀន"
  :showDefaultFooter="false"
  width="600px"
  @close="closeUpdateAdmission"
>
  <AppForm
    :model="formupdateadmission"
    :show-actions="true"
    @submit="submitUpdateAdmission"
    submitText="រក្សាទុក"
  >
    <el-row :gutter="20">
      <el-col :span="8">
        <AppSelect
          v-model="selectacademicforupdateadmission"
          :options="academicOptions"
          placeholder="ឆ្នាំសិក្សា"
          label="ឆ្នាំសិក្សា"
          clearable
          @change="onAcademicChange"
        />
      </el-col>
      <el-col :span="8">
        <AppSelect
          v-model="selectgenerationforupdateadmission"
          :options="generationOptions"
          :disabled="!selectacademicforupdateadmission"
          placeholder="ជំនាន់"
          label="ជំនាន់"
          clearable
          @change="onGenerationChange"
        />
      </el-col>
      <el-col :span="8">
        <AppSelect
          v-model="formupdateadmission.term_id"
          :options="termOptions"
          :disabled="!selectgenerationforupdateadmission"
          placeholder="វគ្គ"
          label="វគ្គ"
          clearable
        />
      </el-col>
    </el-row>
    <el-row>
      <el-col>
        <AppSelect
          v-model="formupdateadmission.state"
          :options="admissionStateOptions"
          placeholder="ស្ថានភាពពាក្យ"
          label="ស្ថានភាពពាក្យ"
          clearable
        />
      </el-col>
    </el-row>
  </AppForm>
</AppDialog>
  </div>
</template>

<style scoped>
.admission-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
