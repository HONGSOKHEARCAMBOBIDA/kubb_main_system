<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getAcademics, createAcademic, updateAcademic, toggleAcademic } from '../../services/academic.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import { debounce } from 'lodash-es'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()

const academics = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusOptions = [
  { label: 'សកម្ម', value: 1 },
  { label: 'អសកម្ម', value: 0 },
]

const filters = reactive({ name: '', active: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះ' },
  { prop: 'start_date', label: 'ថ្ងៃចាប់ផ្តើម', width: 130 },
  { prop: 'end_date', label: 'ថ្ងៃបញ្ចប់', width: 130 },
  { prop: 'description', label: 'ការពិពណ៌នា' },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  code: '',
  name: '',
  start_date: '',
  end_date: '',
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  start_date: [{ required: true, message: 'សូមជ្រើសរើសថ្ងៃចាប់ផ្តើម', trigger: 'change' }],
  // end_date: [{ required: true, message: 'សូមជ្រើសរើសថ្ងៃបញ្ចប់', trigger: 'change' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchAcademics() {
  loading.value = true
  try {
    const res = await getAcademics({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      active: filters.active,
    })
    academics.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchAcademics()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតឆ្នាំសិក្សាថ្មី'
  form.code = ''
  form.name = ''
  form.start_date = ''
  form.end_date = ''
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.uuid
  dialogTitle.value = 'កែប្រែឆ្នាំសិក្សា'
  form.code = row.code
  form.name = row.name
  form.start_date = row.start_date
  form.end_date = row.end_date
  form.description = row.description
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  const valid = await formRef.value?.validate?.().catch(() => false)
  if (valid === false) return

  submitting.value = true
  try {
    const payload = {
      code: form.code,
      name: form.name,
      start_date: form.start_date,
      end_date: form.end_date || null,
      description: form.description,
    }
    if (editingId.value) {
      await updateAcademic(editingId.value, payload)
    } else {
      await createAcademic(payload)
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchAcademics()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleAcademic(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchAcademics()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

onMounted(() => {
  fetchAcademics()
})
</script>

<template>
  <div class="academic-page">
    <AppFilterBar
    :fields="[
      {slot: 'search',span: 10},
      {slot: 'active',span: 5},
      {slot: 'create',span: 5}
    ]"
    >
    <template #search>
        <AppInput
        v-model="filters.name"
        placeholder="ស្វែងរកតាមឈ្មោះ"
        clearable
        @input="debouncedFetch"
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
      show-index
      :data="academics"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchAcademics"
    >
      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="default" @click="openEdit(row)" />
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
        <AppInput v-model="form.code" placeholder="បញ្ចូលលេខកូដ" clearable prop="code" label="លេខកូដ" />
        <AppInput v-model="form.name" placeholder="បញ្ចូលឈ្មោះ" clearable prop="name" label="ឈ្មោះ" />
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
.academic-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>