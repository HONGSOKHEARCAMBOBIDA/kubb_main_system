<script setup>
import { ref, reactive, onMounted } from 'vue'
// import { getprogrammes } from '../../services/programmes.service.js'
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

const programmes = ref([])
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
  { prop: 'name', label: 'ឈ្មោះ' },
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
  name: '',
  description: '',
})

const rules = {
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  description: [{ required: true, message: 'សូមបញ្ចូលការពិពណ៌នា', trigger: 'blur' }],
}

async function fetchProgrammes() {
  loading.value = true
  try {
    const res = await getprogrammes({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      active: filters.active,
    })
    programmes.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchProgrammes()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតកម្មវិធីថ្មី'
  form.name = ''
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែកម្មវិធី'
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
    fetchProgrammes()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    // TODO: call toggle-status service with row.id
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchProgrammes()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

onMounted(() => {
  fetchProgrammes()
})
</script>

<template>
  <div class="category-page">
    <AppFilterBar>
      <AppInput
        v-model="filters.name"
        placeholder="ស្វែងរកតាមឈ្មោះ"
        clearable
        @input="debouncedFetch"
      />
      <AppSelect
        v-model="filters.active"
        :options="statusOptions"
        placeholder="ស្ថានភាព"
        clearable
        @change="debouncedFetch"
      />
      <AppButton type="primary" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="programmes"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchProgrammes"
    >
      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="warning" @click="openEdit(row)" />
        </el-tooltip>
        <el-tooltip content="បិទ/បេីក" placement="top">
          <AppButton
            :icon="row.active ? 'CircleClose' : 'CircleCheck'"
            circle
            size="small"
            type="danger"
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
        <AppInput v-model="form.name" placeholder="បញ្ចូលឈ្មោះ" clearable prop="name" label="ឈ្មោះ" />
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
.category-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>