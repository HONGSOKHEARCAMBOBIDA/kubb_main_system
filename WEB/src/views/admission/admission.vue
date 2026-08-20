<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getAdmission } from '../../services/admission.service'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppButton from '../../components/button/AppButton.vue'

const notify = useNotification()

const admissions = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = reactive({
  student_id: '',
  student_name: '',
})

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
  { prop: 'term_name', label: 'វគ្គ', minwidth: 140 },
  { prop: 'generation_name', label: 'ជំនាន់', slot: 'generation', minwidth: 130 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', slot: 'academic', minwidth: 130 },
  { slot: 'major_name', label: 'ជំនាញ', minwidth: 150 },
  { prop: 'programme_name', label: 'កម្រិត', minwidth: 150 },
  { slot: 'group', label: 'ក្រុមបញ្ចុះតម្លៃ', minwidth: 150 },
  { prop: 'date', label: 'ថ្ងៃដាក់ពាក្យ', minwidth: 120 },
  { label: 'ស្ថានភាពពាក្យសុំ', slot: 'state', minwidth: 130 },
  { label: 'ចំនួនដាក់ពាក្យ', slot: 'enrollCount', minwidth: 130 },
  { label: 'ស្ថានភាព', slot: 'isActive', minwidth: 100 },
]

const columnenrollments = [
  { slot: 'schoolarship', label: 'អាហារូបករណ៍ទទួលបាន', minwidth: 160 },
  { slot: 'fee_interval', label: 'សុំបង់ប្រាក់ជា', minwidth: 130 },
  { slot: 'description', label: 'ផ្សេងៗ' },
]

const columnstudentterms = [
  { slot: 'semester', label: 'ឆមាស', minwidth: 160 },
  {prop:'study_year_id',label:'ឆ្នាំទី',minwidth: 160},
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 130 },
  { prop: 'status', label: 'ស្ថានភាព', slot: 'stStatus', minwidth: 120 },
  { label: 'សកម្ម', slot: 'stActive', minwidth: 90 },
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
  { prop: 'sequence_no', label: 'លេីកទី', minwidth: 70 },
  { prop: 'due_date', label: 'ថ្ងៃត្រូវបង់', minwidth: 110 },
  { label: 'ចំនួនត្រូវបង់', slot: 'instAmount', minwidth: 100 },
  { label: 'ស្ថានភាព', slot: 'instStatus', minwidth: 110 },
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
    total.value = res.data.total || 0
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
          <div>{{ row.student_name_kh }} (<el-text type="warning">{{ row.student_gender }}</el-text>)</div>
          <div>{{ row.student_name_en }}</div>
        </el-text>
      </template>

      <template #generation="{ row }">
        <el-text tag="b" style="color: darkcyan;">{{ row.generation_name || '-' }}</el-text>
      </template>

      <template #academic="{ row }">
        <el-text>{{ row.academic_name || '-' }}</el-text>
      </template>

      <template #state="{ row }">
        <el-text>{{ labelOf(stateLabels, row.state) }}</el-text>
      </template>

      <template #enrollCount="{ row }">
        <el-text tag="b" style="color: dodgerblue;">
          {{ (row.enrollment || []).length }} ដង
        </el-text>
      </template>

      <template #major_name="{ row }">
        <el-text>
          <div>{{ row.major_name }}</div>
          <div><el-text type="primary">{{ row.yearly_fee }}$ /Year</el-text></div>
        </el-text>
      </template>

      <template #group="{ row }">
        <el-text tag="b" style="color: crimson">
          {{
            row.discount_type === 'percentage'
              ? `${row.discount_percentage}%`
              : `${formatMoney(row.discount_amount)}$`
          }}
        </el-text>
      </template>

      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <!-- Level 1: enrollments -->
      <template #expand="{ row }">
        <el-divider content-position="left">ព័ត៌មានដាក់ពាក្យ</el-divider>
        <TableCustom
          expandable
          :data="row.enrollment"
          :columns="columnenrollments"
          :show-pagination="false"
        >
          <template #schoolarship="{ row }">
            <el-text v-if="row.schoolarship_id" tag="b" style="color: crimson">
              {{ row.schoolarship_name }} —
              {{
                row.schoolarship_discount_type === 'percentage'
                  ? `ទទួលការបញ្ចុះតម្លៃ${row.schoolarship_discount_percentage}%`
                  : `ទទួលការបញ្ចុះតម្លៃ${(row.schoolarship_discount_amount)}$`
              }}
            </el-text>
            <el-text v-else type="info">គ្មាន</el-text>
          </template>

          <template #fee_interval="{ row }">
            <el-text tag="b" type="success">{{ labelOf(feeIntervalLabels, row.fee_interval) }}</el-text>
          </template>

          <template #description="{ row }">
            <el-text>{{ row.description || '-' }}</el-text>
          </template>

          <!-- Level 2: student terms -->
          <template #expand="{ row }">
      <el-divider content-position="left">កំពុងសិក្សា</el-divider>
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
                <el-text tag="b" :type="row.status == 'PENDING' ? 'success' : 'info'">{{ row.status }}</el-text>
              </template>

              <template #stActive="{ row }">
                <el-tag :type="row.active ? 'success' : 'danger'" size="small">
                  {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
                </el-tag>
              </template>

              <!-- Level 3: fees -->
              <template #expand="{ row }">
                  <el-divider content-position="left">ថ្លៃសិក្សា</el-divider>
                <TableCustom
                  expandable
                  :data="row.fee"
                  :columns="columnfees"
                  :show-pagination="false"
                >
                  <template #feeAmount="{ row }">
                    <el-text tag="b" style="color: black;">
  {{ formatMoney(row.amount) }}$
                    </el-text>
                  
                  </template>
                  <template #feeDiscount="{ row }">
                    <el-text tag="b" style="color: crimson">
 {{ formatMoney(row.discount) }}$
                    </el-text>
                   
                  
                  </template>
                  <template #feeTotal="{ row }">
                    <el-text tag="b" type="primary">{{ formatMoney(row.total) }}$</el-text>
                  </template>
                  <template #feeActive="{ row }">
                    <el-tag :type="row.active ? 'success' : 'danger'" size="small">
                      {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
                    </el-tag>
                  </template>

                  <!-- Level 4: invoices + installments -->
                  <template #expand="{ row }">
                        <div>
                      <el-divider content-position="left">ការបង់រំលស់</el-divider>
                      <TableCustom
                        expandable
                        :data="row.installment"
                        :columns="columninstallments"
                        :show-pagination="false"
                      >
                        <template #instAmount="{ row }">{{ formatMoney(row.amount) }}$</template>
                        <template #instStatus="{ row }">
                          {{ labelOf(installmentStatusLabels, row.status) }}
                        </template>
                        <template #actions>
                          <AppButton type="success">
                            បង់ប្រាក់
                          </AppButton>
                        </template>
                      </TableCustom>
                    </div>
                    <!-- <div class="mb-3">
                      <el-text tag="b" size="small">វិក័យបត្រ</el-text>
                      <TableCustom
                        expandable
                        :data="row.invoice"
                        :columns="columninvoices"
                        :show-pagination="false"
                      >
                        <template #invGrantTotal="{ row }">
                          <el-text tag="b" type="primary">{{ formatMoney(row.grant_total) }}$</el-text>
                        </template>
                        <template #invTax="{ row }">{{ formatMoney(row.tax) }}$</template>

                       
                        <template #expand="{ row }">
                          <TableCustom
                            :data="row.payment"
                            :columns="columnpayments"
                            :show-pagination="false"
                          >
                            <template #payAmount="{ row }">{{ formatMoney(row.amount) }}$</template>
                          </TableCustom>
                        </template>
                      </TableCustom>
                    </div> -->


                  </template>
                </TableCustom>
              </template>
            </TableCustom>
          </template>
        </TableCustom>
      </template>
    </TableCustom>
  </div>
</template>

<style scoped>
.admission-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>