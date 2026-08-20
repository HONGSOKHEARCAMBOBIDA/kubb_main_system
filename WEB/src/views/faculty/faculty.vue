<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import {
  getFaculty,
  getFacultyByProgrammes,
  createFaculty,
  updateFaculty,
  toggleFaculty,
} from '../../services/faculty.service'
import { getprogrammes } from '../../services/programmes.service.js'
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
const faculty = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusOptions = [
  { label: 'សកម្ម', value: 1 },
  { label: 'អសកម្ម', value: 0 },
]

const programmesOptions = ref([])
const filters = reactive({ programme_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 120 },
  { prop: 'name', label: 'ឈ្មោះមហាវិទ្យាល័យ' },
  { prop: 'programme_name', label: 'កម្រិត', width: 150,slot:'programme_name' },
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
  programme_id: null,
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  programme_id: [{ required: true, message: 'សូមជ្រើសរើសកម្រិត', trigger: 'change' }],
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

async function fetchFaculty() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.programme_id) {
      params.programme_id = filters.programme_id
    }

    const res = await getFaculty(params)
    faculty.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchFaculty()
}, 400)

function openCreate() {
  editingUuid.value = null
  dialogTitle.value = 'បង្កើតមហាវិទ្យាល័យថ្មី'
  form.code = ''
  form.name = ''
  form.programme_id = null
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  // backend looks rows up by uuid, not the numeric id
  editingUuid.value = row.uuid
  dialogTitle.value = 'កែប្រែមហាវិទ្យាល័យ'
  form.code = row.code
  form.name = row.name
  form.programme_id = row.programme_id ?? null
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
      // end_date only makes sense on update, per backend logic
      await updateFaculty(editingUuid.value, {
        code: form.code,
        name: form.name,
        programme_id: form.programme_id,
        description: form.description,
      })
    } else {
      await createFaculty({
        code: form.code,
        name: form.name,
        programme_id: form.programme_id,
        description: form.description,
      })
    }
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchFaculty()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleFaculty(row.uuid)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchFaculty()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

watch(
  () => filters.programme_id,
  () => {
    page.value = 1
    fetchFaculty()
  }
)

onMounted(() => {
  fetchProgrammeOptions()
  fetchFaculty()
})
</script>

<template>
  <div class="generation-page">
    <AppFilterBar
      :fields="[
        { slot: 'program', span: 5 },
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
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="faculty"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchFaculty"
    >
      <template #programme_name="{ row }">
        <el-text tag="b" style="color: darkcyan;">
          {{ row.programme_name || '-' }}
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
          label="ឈ្មោះមហាវិទ្យាល័យ"
        />
        <AppSelect
          v-model="form.programme_id"
          :options="programmesOptions"
          placeholder="ជ្រើសរើសកម្រិតសិក្សា"
          clearable
          prop="programme_id"
          label="កម្រិតសិក្សា"
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
.generation-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>
