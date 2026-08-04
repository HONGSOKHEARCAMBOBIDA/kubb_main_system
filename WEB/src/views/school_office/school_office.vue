<script setup>
import { ref, reactive, onMounted } from 'vue'
// import {getSchoolOffices} from '../../services/school_office.service.js'
import { getSchoolOffices } from '../../services/school_office.service.js'
import { getFloors } from '../../services/floor.service.js'
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

const schoolOffices = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// floor options for the select filter/form (loaded once)
const floorOptions = ref([])

const filters = reactive({ name: '', floor_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 110 },
  { prop: 'name', label: 'ឈ្មោះការិយាល័យ' },
  { prop: 'floor', label: 'ឈ្មោះជាន់', slot: 'floorName', width: 140 },
  { prop: 'building', label: 'ឈ្មោះអគារ', slot: 'buildingName', width: 160 },
  { prop: 'campuse', label: 'ឈ្មោះសាខា', slot: 'campusName', width: 160 },
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
  floor_id: null,
  code: '',
  name: '',
  address: '',
  description: '',
})

const rules = {
  floor_id: [{ required: true, message: 'សូមជ្រើសរើសជាន់', trigger: 'change' }],
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះការិយាល័យ', trigger: 'blur' }],
  address: [{ required: true, message: 'សូមបញ្ចូលអាសយដ្ឋាន', trigger: 'blur' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchFloorOptions() {
  try {
    const res = await getFloors({ page: 1, pageSize: 100 })
    floorOptions.value = (res.data.data || []).map((f) => {
      const buildingName = f.building?.name
      const campusName = f.building?.campuse?.name
      const parts = [buildingName, campusName].filter(Boolean)
      return {
        label: parts.length ? `${f.name} (${parts.join(' - ')})` : f.name,
        value: f.id,
      }
    })
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load floors')
  }
}

async function fetchSchoolOffices() {
  loading.value = true
  try {
    const res = await getSchoolOffices({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      floor_id: filters.floor_id,
    })
    schoolOffices.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchSchoolOffices()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតការិយាល័យថ្មី'
  form.floor_id = null
  form.code = ''
  form.name = ''
  form.address = ''
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែការិយាល័យ'
  form.floor_id = row.floor?.id ?? row.floor_id ?? null
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
    fetchSchoolOffices()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchFloorOptions()
  fetchSchoolOffices()
})
</script>

<template>
  <div class="school-office-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 10 },
        { slot: 'floor', span: 8 },
      ]"
      :action-span="6"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះការិយាល័យ"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <template #floor>
        <AppSelect
          v-model="filters.floor_id"
          :options="floorOptions"
          placeholder="ជាន់"
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
      :data="schoolOffices"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchSchoolOffices"
    >
      <template #floorName="{ row }">
        <el-text tag="mark">
          {{ row.floor?.name || '-' }}
        </el-text>
      </template>

      <template #buildingName="{ row }">
        <el-text tag="mark">
          {{ row.floor?.building?.name || '-' }}
        </el-text>
      </template>

      <template #campusName="{ row }">
        <el-text tag="mark">
          {{ row.floor?.building?.campuse?.name || '-' }}
        </el-text>
      </template>

      <template #schoolName="{ row }">
        <el-text tag="mark">
          {{ row.floor?.building?.campuse?.school?.name || '-' }}
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
          v-model="form.floor_id"
          :options="floorOptions"
          placeholder="ជ្រើសរើសជាន់"
          clearable
          prop="floor_id"
          label="ជាន់"
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
          placeholder="បញ្ចូលឈ្មោះការិយាល័យ"
          clearable
          prop="name"
          label="ឈ្មោះការិយាល័យ"
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
.school-office-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>