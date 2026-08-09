<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import {
  getDepartment,
  createDepartment,
  updateDepartment,
  toggleDepartment,
} from '../../services/department.service.js'
import { getprogrammes } from '../../services/programmes.service.js'
import { getFacultyByProgrammes } from '../../services/faculty.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()
const department = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const formProgramID = ref(null)
const programmesOptions = ref([])
const facultyOptions = ref([])
const filters = reactive({ faculty_id: null, programme_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះដេប៉ាតេម៉ង' },
  { slot: 'programme_name', label: 'កម្រិត', width: 150 },
  { prop: 'faculty_name', label: 'មហាវិទ្យាល័យ', width: 250, slot: 'faculty_name' },
  { prop: 'description', label: 'ការពិពណ៌នា' },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingUuid = ref(null)
const isEditing = computed(() => !!editingUuid.value)

const form = reactive({
  code: '',
  name: '',
  faculty_id: null,
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  faculty_id: [{ required: true, message: 'សូមជ្រើសរើសមហាវិទ្យាល័យ', trigger: 'change' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchProgrammeOptions() {
  try {
    const res = await getprogrammes()
    programmesOptions.value = (res.data.data || []).map((a) => ({
      label: a.name,
      value: a.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
  }
}

async function loadFacultyOption(programmeID) {
  if (!programmeID) return []
  try {
    const res = await getFacultyByProgrammes(programmeID)
    return (res.data.data || []).map((f) => ({
      label: f.name,
      value: f.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load faculties')
    return []
  }
}

async function fetchDepartment() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.faculty_id) {
      params.faculty_id = filters.faculty_id
    }

    const res = await getDepartment(params)
    department.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតដេប៉ាតេម៉ងថ្មី'
  form.code = ''
  form.name = ''
  form.faculty_id = null
  form.description = ''
  formProgramID.value = null
  facultyOptions.value = []
  dialogVisible.value = true
}

async function openEdit(row) {
  // backend looks rows up by uuid, not the numeric id
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែដេប៉ាតេម៉ង'
  form.code = row.code
  form.name = row.name
  form.faculty_id = row.faculty_id ?? null
  form.description = row.description

  // We don't know the row's programme_id from this response, so we can't
  // preselect the programme dropdown. Seed facultyOptions with just the
  // department's current faculty so the select shows the right label.
  formProgramID.value = row.programme_id ?? null
  facultyOptions.value = await loadFacultyOption(formProgramID.value)
  form.faculty_id = row.faculty_id ?? null
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (isEditing.value) {
      await updateDepartment(editingUuid.value, {
        code: form.code,
        name: form.name,
        faculty_id: form.faculty_id,
        description: form.description,
      })
    } else {
      await createDepartment({
        code: form.code,
        name: form.name,
        faculty_id: form.faculty_id,
        description: form.description,
      })
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchDepartment()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleDepartment(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchDepartment()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

watch(
  () => filters.programme_id,
  async (newVal) => {
    filters.faculty_id = null
    facultyOptions.value = await loadFacultyOption(newVal)
    page.value = 1
    fetchDepartment()
  }
)

watch(
  () => filters.faculty_id,
  () => {
    page.value = 1
    fetchDepartment()
  }
)

watch(formProgramID, async (newVal) => {
  form.faculty_id = null
  facultyOptions.value = await loadFacultyOption(newVal)
})

onMounted(() => {
  fetchProgrammeOptions()
  fetchDepartment()
})
</script>

<template>
  <div class="department-page">
    <AppFilterBar
      :fields="[
        { slot: 'program', span: 5 },
        { slot: 'faculty', span: 5 },
        { slot: 'create', span: 5 },
      ]"
      :action-span="3"
    >
      <template #program>
        <AppSelect
          v-model="filters.programme_id"
          :options="programmesOptions"
          placeholder="កម្រិតសិក្សា"
          clearable
        />
      </template>
      <template #faculty>
        <AppSelect
          v-model="filters.faculty_id"
          :options="facultyOptions"
          placeholder="មហាវិទ្យាល័យ"
          clearable
        />
      </template>
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="department"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchDepartment"
    >
      <template #programme_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.programme_name }}
        </el-text>
      </template>
      <template #faculty_name="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.faculty_code }} - {{ row.faculty_name }}
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
          label="ឈ្មោះដេប៉ាតេម៉ង"
        />
        <AppSelect
          v-model="formProgramID"
          :options="programmesOptions"
          placeholder="ជ្រើសរើសកម្រិតសិក្សា"
          clearable
          prop="programme_id"
          label="កម្រិតសិក្សា"
        />
        <AppSelect
          v-model="form.faculty_id"
          :options="facultyOptions"
          placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
          clearable
          prop="faculty_id"
          label="មហាវិទ្យាល័យ"
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
.department-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>
