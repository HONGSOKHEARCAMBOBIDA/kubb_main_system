<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { getTerm, createTerm, updateTerm, toggleTerm } from '../../services/term.service'
import { getAcademics } from '../../services/academic.service'
import { getGenerationByAcademic } from '../../services/generation.service'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()

const terms = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusOptions = [
  { label: 'សកម្ម', value: 1 },
  { label: 'អសកម្ម', value: 0 },
]

const academicOptions = ref([])
const filterGenerationOptions = ref([])
const filters = reactive({ academic_id: null, generation_id: null, active: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះវគ្គ' },
  { prop: 'generation_name', label: 'ឈ្មោះជំនាន់', slot: 'generationName', width: 150 },
  { prop: 'academic_name', label: 'ឈ្មោះឆ្នាំសិក្សា', slot: 'academicName', width: 150 },
  { prop: 'start_date', label: 'ថ្ងៃចាប់ផ្តើម', width: 130 },
  {  label: 'ថ្ងៃបញ្ចប់', width: 130,slot: 'enddate' },
  { prop: 'description', label: 'ការពិពណ៌នា' },
  {label:'ជំនាញកំពុងបើក',slot: 'majors',width: 200},
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const formAcademicId = ref(null)
const formGenerationOptions = ref([])

const form = reactive({
  code: '',
  name: '',
  generation_id: null,
  start_date: '',
  end_date: '',
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  generation_id: [{ required: true, message: 'សូមជ្រើសរើសជំនាន់', trigger: 'change' }],
  start_date: [{ required: true, message: 'សូមជ្រើសរើសថ្ងៃចាប់ផ្តើម', trigger: 'change' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchAcademicOptions() {
  try {
    const res = await getAcademics()
    academicOptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
  }
}

// uses the dedicated by-academic endpoint, not the paginated list endpoint
async function loadGenerationOptions(academicId) {
  if (!academicId) return []
  try {
    const res = await getGenerationByAcademic(academicId)
    return (res.data.data || []).map((g) => ({
      label: g.name,
      value: g.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load generations')
    return []
  }
}

async function fetchTerm() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.academic_id) params.academic_id = filters.academic_id
    if (filters.generation_id) params.generation_id = filters.generation_id
    if (filters.active !== null && filters.active !== undefined) params.active = filters.active

    const res = await getTerm(params)
    terms.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

// ---- filter bar cascading via watchers (robust to clear/×) ----
watch(
  () => filters.academic_id,
  async (newVal) => {
    filters.generation_id = null
    filterGenerationOptions.value = await loadGenerationOptions(newVal)
    page.value = 1
    fetchTerm()
  }
)

watch(
  () => filters.generation_id,
  () => {
    page.value = 1
    fetchTerm()
  }
)

watch(
  () => filters.active,
  () => {
    page.value = 1
    fetchTerm()
  }
)

// ---- form cascading ----
watch(formAcademicId, async (newVal) => {
  form.generation_id = null
  formGenerationOptions.value = await loadGenerationOptions(newVal)
})

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតវគ្គថ្មី'
  form.code = ''
  form.name = ''
  form.generation_id = null
  form.start_date = ''
  form.end_date = ''
  form.description = ''
  formAcademicId.value = null
  formGenerationOptions.value = []
  dialogVisible.value = true
}

async function openEdit(row) {
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែវគ្គ'
  form.code = row.code
  form.name = row.name
  form.start_date = row.start_date
  form.end_date = row.end_date
  form.description = row.description

  formAcademicId.value = row.academic_id ?? null
  formGenerationOptions.value = await loadGenerationOptions(formAcademicId.value)
  form.generation_id = row.generation_id ?? null

  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (isEditing.value) {
      await updateTerm(editingUuid.value, {
        code: form.code,
        name: form.name,
        generation_id: form.generation_id,
        start_date: form.start_date,
        end_date: form.end_date || null,
        description: form.description,
      })
    } else {
      await createTerm({
        code: form.code,
        name: form.name,
        generation_id: form.generation_id,
        start_date: form.start_date,
        end_date: form.end_date,
        description: form.description,
      })
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchTerm()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleTerm(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchTerm()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

onMounted(() => {
  fetchAcademicOptions()
  fetchTerm()
})
</script>

<template>
  <div class="term-page">
    <AppFilterBar
      :fields="[
        { slot: 'academic', span: 5 },
        { slot: 'generation', span: 5 },
        { slot: 'active', span: 5 },
        { slot: 'create', span: 5 },
      ]"
      :action-span="3"
    >
      <template #academic>
        <AppSelect
          v-model="filters.academic_id"
          :options="academicOptions"
          placeholder="ឆ្នាំសិក្សា"
          clearable
          @change="onFilterAcademicChange"
        />
      </template>
      <template #generation>
        <AppSelect
          v-model="filters.generation_id"
          :options="filterGenerationOptions"
          placeholder="ជំនាន់"
          clearable
          :disabled="!filters.academic_id"
          @change="onFilterGenerationChange"
        />
      </template>
      <template #active>
        <AppSelect
          v-model="filters.active"
          :options="statusOptions"
          placeholder="ស្ថានភាព"
          clearable
          @change="debouncedFetch"
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      expandable
      :data="terms"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchTerm"
    >
      <template #generationName="{ row }">
        <el-text tag="b" style="color: darkcyan;">
          {{ row.generation_name || '-' }}
        </el-text>
      </template>

      <template #academicName="{ row }">
        <el-text>{{ row.academic_name || '-' }}</el-text>
      </template>

      <template #enddate="{row}">
         <el-text tag="b" style="color: crimson;">{{ row.end_date || '-' }}</el-text>
      </template>

      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>
      <template #majors="{row}">
        <el-text tag="b" style="color: dodgerblue;">
            {{ row.majors.length }} ជំនាញ
        </el-text>
      </template>

<template #expand="{ row }">
  <div class="p-4 bg-gray-50">
    <div class="flex items-center justify-between mb-3">
      <div>
        <el-text tag="b">ជំនាញ</el-text>
        <el-text type="info" class="ml-2">
          {{ row.majors?.length || 0 }} កំពុងបើក
        </el-text>
      </div>
    </div>

    <div
      v-if="row.majors?.length"
      class="grid grid-cols-1 md:grid-cols-2 gap-3"
    >
      <div
        v-for="major in row.majors"
        :key="major.id"
        class="rounded-lg border bg-white p-3"
      >
        <div class="flex items-center gap-2">
          <el-tag size="small">
            {{ major.code }}
          </el-tag>

          <el-text tag="b">
            {{ major.name }}
          </el-text>
        </div>

        <div class="mt-2 text-sm text-gray-500">
          ដេប៉ាតម៉ង.
          {{ major.department_name }}
        </div>

        <div class="text-xs text-gray-400 mt-1">
          {{ major.faculty_name }}
          ·
          {{ major.programme_name }}
        </div>
      </div>
    </div>

    <el-empty
      v-else
      description="មិនទាន់មានជំនាញ"
      :image-size="60"
    />
  </div>
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
          label="ឈ្មោះវគ្គ"
        />
        <AppSelect
          v-model="formAcademicId"
          :options="academicOptions"
          placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
          clearable
          label="ឆ្នាំសិក្សា"
          @change="onFormAcademicChange"
        />
        <AppSelect
          v-model="form.generation_id"
          :options="formGenerationOptions"
          placeholder="ជ្រើសរើសជំនាន់"
          clearable
          prop="generation_id"
          label="ជំនាន់"
          :disabled="!formAcademicId"
        />
        <AppInput
          v-model="form.start_date"
          type="date"
          placeholder="ជ្រើសរើសថ្ងៃចាប់ផ្តើម"
          clearable
          prop="start_date"
          label="ថ្ងៃចាប់ផ្តើម"
        />
        <AppInput
          v-model="form.end_date"
          type="date"
          placeholder="ជ្រើសរើសថ្ងៃបញ្ចប់"
          clearable
          prop="end_date"
          label="ថ្ងៃបញ្ចប់"
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
.term-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>