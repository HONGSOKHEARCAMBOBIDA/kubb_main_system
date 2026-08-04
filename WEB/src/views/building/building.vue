<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBuildings } from '../../services/building.service.js'
import { getCampuses } from '../../services/campuse.service.js'
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

const buildings = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// campus options for the select filter/form (loaded once)
const campusOptions = ref([])

const filters = reactive({ name: '', campus_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 110 },
  { prop: 'name', label: 'ឈ្មោះអគារ' },
  { prop: 'campuse', label: 'ឈ្មោះសាលា', slot: 'campusName', width: 160 },
  { prop: 'school', label: 'សាខាគោល', slot: 'schoolName', width: 180 },
  { prop: 'address', label: 'អាសយដ្ឋាន' },
  { prop: 'description', label: 'ការពិពណ៌នា' },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  campus_id: null,
  code: '',
  name: '',
  address: '',
  description: '',
})

const rules = {
  campus_id: [{ required: true, message: 'សូមជ្រើសរើសសាខា', trigger: 'change' }],
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះអគារ', trigger: 'blur' }],
  address: [{ required: true, message: 'សូមបញ្ចូលអាសយដ្ឋាន', trigger: 'blur' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchCampusOptions() {
  try {
    const res = await getCampuses({ page: 1, pageSize: 100 })
    campusOptions.value = (res.data.data || []).map((c) => ({
      label: c.school ? `${c.name} (${c.school.name})` : c.name,
      value: c.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load campuses')
  }
}

async function fetchBuildings() {
  loading.value = true
  try {
    const res = await getBuildings({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      campus_id: filters.campus_id,
    })
    buildings.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchBuildings()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតអគារថ្មី'
  form.campus_id = null
  form.code = ''
  form.name = ''
  form.address = ''
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែអគារ'
  form.campus_id = row.campuse?.id ?? row.campus_id ?? null
  form.code = row.code
  form.name = row.name
  form.address = row.address
  form.description = row.description
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    // TODO: call create/update service depending on editingId.value
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchBuildings()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCampusOptions()
  fetchBuildings()
})
</script>

<template>
  <div class="building-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 10 },
        { slot: 'campus', span: 8 },
      ]"
      :action-span="6"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះអគារ"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <template #campus>
        <AppSelect
          v-model="filters.campus_id"
          :options="campusOptions"
          placeholder="សាខា"
          clearable
          @change="debouncedFetch"
        />
      </template>
      <template #actions>
        <AppButton type="primary" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="buildings"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchBuildings"
    >
      <template #campusName="{ row }">
     <el-text tag="mark">
           {{ row.campuse?.name || '-' }}
     </el-text>
      </template>

      <template #schoolName="{ row }">
       <el-text tag="mark">
         {{ row.campuse?.school?.name || '-' }}
       </el-text>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="warning" @click="openEdit(row)" />
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
        <AppSelect
          v-model="form.campus_id"
          :options="campusOptions"
          placeholder="ជ្រើសរើសសាខា"
          clearable
          prop="campus_id"
          label="សាខា"
        />
        <AppInput
          v-model="form.code"
          placeholder="បញ្ចូលលេខកូដ"
          clearable
          prop="code"
          label="លេខកូដ"
        />
        <AppInput
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះអគារ"
          clearable
          prop="name"
          label="ឈ្មោះអគារ"
        />
        <AppInput
          v-model="form.address"
          placeholder="បញ្ចូលអាសយដ្ឋាន"
          clearable
          prop="address"
          label="អាសយដ្ឋាន"
          type="textarea"
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
.building-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>