<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getFloors } from '../../services/floor.service.js'
import { getBuildings } from '../../services/building.service.js'
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

const floors = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// building options for the select filter/form (loaded once)
const buildingOptions = ref([])

const filters = reactive({ name: '', building_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 110 },
  { prop: 'name', label: 'ឈ្មោះជាន់' },
  { prop: 'building', label: 'ឈ្មោះអគារ', slot: 'buildingName', width: 160 },
  { prop: 'campuse', label: 'ឈ្មោះសាខា', slot: 'campusName', width: 160 },
  { prop: 'school', label: 'សាខាគោល', slot: 'schoolName', width: 180 },
  { prop: 'description', label: 'ការពិពណ៌នា' },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  building_id: null,
  code: '',
  name: '',
  description: '',
})

const rules = {
  building_id: [{ required: true, message: 'សូមជ្រើសរើសអគារ', trigger: 'change' }],
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះជាន់', trigger: 'blur' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchBuildingOptions() {
  try {
    const res = await getBuildings({ page: 1, pageSize: 100 })
    buildingOptions.value = (res.data.data || []).map((b) => ({
      label: b.campuse?.school
        ? `${b.name} (${b.campuse.name} - ${b.campuse.school.name})`
        : b.campuse
          ? `${b.name} (${b.campuse.name})`
          : b.name,
      value: b.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load buildings')
  }
}

async function fetchFloors() {
  loading.value = true
  try {
    const res = await getFloors({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      building_id: filters.building_id,
    })
    floors.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchFloors()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតជាន់ថ្មី'
  form.building_id = null
  form.code = ''
  form.name = ''
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែជាន់'
  form.building_id = row.building?.id ?? row.building_id ?? null
  form.code = row.code
  form.name = row.name
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
    fetchFloors()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchBuildingOptions()
  fetchFloors()
})
</script>

<template>
  <div class="floor-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 10 },
        { slot: 'building', span: 8 },
      ]"
      :action-span="6"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះជាន់"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <template #building>
        <AppSelect
          v-model="filters.building_id"
          :options="buildingOptions"
          placeholder="អគារ"
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
      :data="floors"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchFloors"
    >
      <template #buildingName="{ row }">
        <el-text tag="mark">
          {{ row.building?.name || '-' }}
        </el-text>
      </template>

      <template #campusName="{ row }">
        <el-text tag="mark">
          {{ row.building?.campuse?.name || '-' }}
        </el-text>
      </template>

      <template #schoolName="{ row }">
        <el-text tag="mark">
          {{ row.building?.campuse?.school?.name || '-' }}
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
          v-model="form.building_id"
          :options="buildingOptions"
          placeholder="ជ្រើសរើសអគារ"
          clearable
          prop="building_id"
          label="អគារ"
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
          placeholder="បញ្ចូលឈ្មោះជាន់"
          clearable
          prop="name"
          label="ឈ្មោះជាន់"
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
.floor-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>