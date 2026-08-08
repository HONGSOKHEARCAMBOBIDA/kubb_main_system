<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { getSemester,createSemester,updateSemester,toggleSemester } from '../../services/semester.service.js'
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

const generations = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusOptions = [
  { label: 'សកម្ម', value: 1 },
  { label: 'អសកម្ម', value: 0 },
]

const academicOptions = ref([])
const filters = reactive({ academic_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះឆមាស' },
  // { prop: 'index', label: 'លំដាប់', width: 90 },
  { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', slot: 'academicName', width: 160 },
  { prop: 'start_date', label: 'ថ្ងៃចាប់ផ្តើម', width: 130 },
  { prop: 'end_date', label: 'ថ្ងៃបញ្ចប់', width: 130,slot:'end_date' },
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
  code: '',
  name: '',
  academic_id: null,
  start_date: '',
  end_date: '',
})

// index removed: server computes it on create and ignores it on update
const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  academic_id: [{ required: true, message: 'សូមជ្រើសរើសឆ្នាំសិក្សា', trigger: 'change' }],
  start_date: [{ required: true, message: 'សូមជ្រើសរើសថ្ងៃចាប់ផ្តើម', trigger: 'change' }],
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
async function fetchSemester() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.academic_id) {
      params.academic_id = filters.academic_id
    }

    const res = await getSemester(params)
    generations.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchSemester()
}, 400)

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតជំនាន់ថ្មី'
  form.code = ''
  form.name = ''
  form.academic_id = null
  form.start_date = ''
  form.end_date = ''
  dialogVisible.value = true
}

function openEdit(row) {
  // backend looks rows up by uuid, not the numeric id
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែជំនាន់'
  form.code = row.code
  form.name = row.name
  form.academic_id = row.academic?.id ?? row.academic_id ?? null
  form.start_date = row.start_date
  form.end_date = row.end_date
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (isEditing.value) {
      // end_date only makes sense on update, per backend logic
      await updateSemester(editingUuid.value, {
        code: form.code,
        name: form.name,
        academic_id: form.academic_id,
        start_date: form.start_date,
        end_date: form.end_date || null,
      })
    } else {
      await createSemester({
        code: form.code,
        name: form.name,
        academic_id: form.academic_id,
        start_date: form.start_date,
        end_date: form.end_date || null,
      })
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchSemester()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleSemester(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchSemester()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

watch(
  () => filters.academic_id,
  () => {
    page.value = 1
    fetchSemester()
  }
)


onMounted(() => {
  fetchAcademicOptions()
  fetchSemester()
})
</script>

<template>
  <div class="generation-page">
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

    <TableCustom show-index :data="generations" :columns="columns" :loading="loading" :total="total"
      v-model:current-page="page" v-model:page-size="pageSize" @page-change="fetchGenerations">
      <template #academicName="{ row }">
        <el-text tag="b" style="color: cornflowerblue;">
          {{ row.academic_name || '-' }}
        </el-text>
      </template>
      <template #end_date="{row}">
        <el-text tag="b" style="color: crimson;">{{ row.end_date }}</el-text>
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
        <AppInput v-model="form.code" placeholder="បញ្ចូលលេខកូដ" clearable prop="code" label="លេខកូដ" />
        <AppInput v-model="form.name" placeholder="បញ្ចូលឈ្មោះ" clearable prop="name" label="ឈ្មោះជំនាន់" />
        <!-- <AppInput
          v-model.number="form.index"
          type="number"
          placeholder="បញ្ចូលលំដាប់"
          clearable
          prop="index"
          label="លំដាប់"
        /> -->
        <AppSelect v-model="form.academic_id" :options="academicOptions" placeholder="ជ្រើសរើសឆ្នាំសិក្សា" clearable
          prop="academic_id" label="ឆ្នាំសិក្សា" />
        <AppInput v-model="form.start_date" type="date" placeholder="ជ្រើសរើសថ្ងៃចាប់ផ្តើម" clearable prop="start_date"
          label="ថ្ងៃចាប់ផ្តើម" />
        <AppInput v-model="form.end_date" type="date" placeholder="ជ្រើសរើសថ្ងៃបញ្ចប់" clearable
          prop="end_date" label="ថ្ងៃបញ្ចប់" />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.generation-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>
