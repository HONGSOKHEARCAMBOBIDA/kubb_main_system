<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import {
  getAcademicshift,
  createAcademicshift,
  updateAcademicshift,
  toggleAcademicshift,
} from '../../services/academic_shift.service.js'

import { getAcademics } from '../../services/academic.service.js'
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

const shifts = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const academicOptions = ref([])
const filters = reactive({ academic_id: null })

const columns = [
  { prop: 'name', label: 'ឈ្មោះវេន' },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', slot: 'academicName', width: 160 },
  { prop: 'description', label: 'ការពិពណ៌នា' },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const form = reactive({
  name: '',
  academic_id: null,
  description: '',
})

const rules = {
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  academic_id: [{ required: true, message: 'សូមជ្រើសរើសឆ្នាំសិក្សា', trigger: 'change' }],
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

async function fetchShifts() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.academic_id) {
      params.academic_id = filters.academic_id
    }

    const res = await getAcademicshift(params)
    shifts.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchShifts()
}, 400)

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតវេនថ្មី'
  form.name = ''
  form.academic_id = null
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  // backend looks rows up by uuid, not the numeric id
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែវេន'
  form.name = row.name
  form.academic_id = row.academic_id ?? null
  form.description = row.description
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (isEditing.value) {
      await updateAcademicshift(editingUuid.value, {
        name: form.name,
        academic_id: form.academic_id,
        description: form.description,
      })
    } else {
      await createAcademicshift({
        name: form.name,
        academic_id: form.academic_id,
        description: form.description,
      })
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchShifts()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleAcademicshift(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchShifts()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

watch(
  () => filters.academic_id,
  () => {
    page.value = 1
    fetchShifts()
  }
)

onMounted(() => {
  fetchAcademicOptions()
  fetchShifts()
})
</script>

<template>
  <div class="academic-shift-page">
    <AppFilterBar :fields="[
      { slot: 'academic', span: 5 },
      { slot: 'create', span: 5 },
    ]" :action-span="3">
      <template #academic>
        <AppSelect v-model="filters.academic_id" :options="academicOptions" placeholder="ឆ្នាំសិក្សា" clearable />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom show-index :data="shifts" :columns="columns" :loading="loading" :total="total"
      v-model:current-page="page" v-model:page-size="pageSize" @page-change="fetchShifts">
      <template #academicName="{ row }">
        <el-text tag="b" style="color: darkcyan;">
          {{ row.academic_name || '-' }}
        </el-text>
      </template>

      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="default" plain @click="openEdit(row)" />
        </el-tooltip>
        <el-tooltip content="បិទ/បេីក" placement="top">
          <AppButton :icon="row.active ? 'CircleClose' : 'CircleCheck'" circle size="small" type="default" plain
            @click="toggleStatus(row)" />
        </el-tooltip>
      </template>
    </TableCustom>

    <AppDialog v-model:visible="dialogVisible" :title="dialogTitle" :showDefaultFooter="false">
      <AppForm ref="formRef" :model="form" :rules="rules" :loading="submitting" :show-actions="true"
        @submit="handleSubmit" @reset="closeDialog" submitText="រក្សាទុក" resetText="ចាកចេញ">
        <AppInput v-model="form.name" placeholder="បញ្ចូលឈ្មោះ" clearable prop="name" label="ឈ្មោះវេន" />
        <AppSelect v-model="form.academic_id" :options="academicOptions" placeholder="ជ្រើសរើសឆ្នាំសិក្សា" clearable
          prop="academic_id" label="ឆ្នាំសិក្សា" />
        <AppInput v-model="form.description" placeholder="បញ្ចូលការពិពណ៌នា" clearable prop="description"
          label="ការពិពណ៌នា" type="textarea" />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.academic-shift-page {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>