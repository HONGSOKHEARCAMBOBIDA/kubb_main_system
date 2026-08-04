<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getCampuses } from '../../services/campuse.service.js'
import { getSchools } from '../../services/school.service.js'
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

const campuses = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// school options for the select filter/form (loaded once)
const schoolOptions = ref([])

const filters = reactive({ name: '', school_id: null })

const columns = [
  { prop: 'code', label: 'លេខកូដ', width: 110 },
  { prop: 'prefix', label: 'បុព្វបទ', width: 90 },
  { prop: 'name', label: 'ឈ្មោះសាខា' },
  { prop: 'school', label: 'សាខាគោល', slot: 'schoolName', width: 180 },
  { prop: 'director', label: 'នាយិកា' },
  { prop: 'school_level', label: 'កម្រិតសិក្សា', width: 130 },
  { prop: 'capital_province', label: 'ខេត្ត/ក្រុង', width: 130 },
  { prop: 'phone', label: 'ទូរស័ព្ទ', width: 130 },
  { prop: 'email', label: 'អ៊ីមែល' },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  school_id: null,
  code: '',
  prefix: '',
  name: '',
  director: '',
  capital_province: '',
  school_level: '',
  phone: '',
  email: '',
  website: '',
  facebook: '',
  address: '',
})

const rules = {
  school_id: [{ required: true, message: 'សូមជ្រើសរើសសាលា', trigger: 'change' }],
  code: [{ required: true, message: 'សូមបញ្ចូលលេខកូដ', trigger: 'blur' }],
  prefix: [{ required: true, message: 'សូមបញ្ចូលបុព្វបទ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះសាខា', trigger: 'blur' }],
  director: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះនាយក', trigger: 'blur' }],
  capital_province: [{ required: true, message: 'សូមបញ្ចូលខេត្ត/ក្រុង', trigger: 'blur' }],
  school_level: [{ required: true, message: 'សូមបញ្ចូលកម្រិតសិក្សា', trigger: 'blur' }],
  phone: [{ required: true, message: 'សូមបញ្ចូលលេខទូរស័ព្ទ', trigger: 'blur' }],
  email: [
    { required: true, message: 'សូមបញ្ចូលអ៊ីមែល', trigger: 'blur' },
    { type: 'email', message: 'អ៊ីមែលមិនត្រឹមត្រូវ', trigger: 'blur' },
  ],
  address: [{ required: true, message: 'សូមបញ្ចូលអាសយដ្ឋាន', trigger: 'blur' }],
}

async function fetchSchoolOptions() {
  try {
    const res = await getSchools({ page: 1, pageSize: 100 })
    schoolOptions.value = (res.data.data || []).map((s) => ({
      label: s.name,
      value: s.id,
    }))
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load schools')
  }
}

async function fetchCampuses() {
  loading.value = true
  try {
    const res = await getCampuses({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      school_id: filters.school_id,
    })
    campuses.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchCampuses()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតសាខាថ្មី'
  form.school_id = null
  form.code = ''
  form.prefix = ''
  form.name = ''
  form.director = ''
  form.capital_province = ''
  form.school_level = ''
  form.phone = ''
  form.email = ''
  form.website = ''
  form.facebook = ''
  form.address = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែសាខា'
  form.school_id = row.school?.id ?? row.school_id ?? null
  form.code = row.code
  form.prefix = row.prefix
  form.name = row.name
  form.director = row.director
  form.capital_province = row.capital_province
  form.school_level = row.school_level
  form.phone = row.phone
  form.email = row.email
  form.website = row.website
  form.facebook = row.facebook
  form.address = row.address
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
    fetchCampuses()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchSchoolOptions()
  fetchCampuses()
})
</script>

<template>
  <div class="campuse-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 10 },
        { slot: 'school', span: 8 },
      ]"
      :action-span="6"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះសាខា"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <template #school>
        <AppSelect
          v-model="filters.school_id"
          :options="schoolOptions"
          placeholder="សាលា"
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
      :data="campuses"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchCampuses"
    >
      <template #schoolName="{ row }">
        <el-text tag="mark">
          {{ row.school?.name || '-' }}
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
          v-model="form.school_id"
          :options="schoolOptions"
          placeholder="ជ្រើសរើសសាលា"
          clearable
          prop="school_id"
          label="សាលា"
        />
        <AppInput
          v-model="form.code"
          placeholder="បញ្ចូលលេខកូដ"
          clearable
          prop="code"
          label="លេខកូដ"
        />
        <AppInput
          v-model="form.prefix"
          placeholder="បញ្ចូលបុព្វបទ"
          clearable
          prop="prefix"
          label="បុព្វបទ"
        />
        <AppInput
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះសាខា"
          clearable
          prop="name"
          label="ឈ្មោះសាខា"
        />
        <AppInput
          v-model="form.director"
          placeholder="បញ្ចូលឈ្មោះនាយក"
          clearable
          prop="director"
          label="នាយក"
        />
        <AppInput
          v-model="form.school_level"
          placeholder="បញ្ចូលកម្រិតសិក្សា"
          clearable
          prop="school_level"
          label="កម្រិតសិក្សា"
        />
        <AppInput
          v-model="form.capital_province"
          placeholder="បញ្ចូលខេត្ត/ក្រុង"
          clearable
          prop="capital_province"
          label="ខេត្ត/ក្រុង"
        />
        <AppInput
          v-model="form.phone"
          placeholder="បញ្ចូលលេខទូរស័ព្ទ"
          clearable
          prop="phone"
          label="ទូរស័ព្ទ"
        />
        <AppInput
          v-model="form.email"
          placeholder="បញ្ចូលអ៊ីមែល"
          clearable
          prop="email"
          label="អ៊ីមែល"
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
          v-model="form.website"
          placeholder="បញ្ចូលគេហទំព័រ (មិនចាំបាច់)"
          clearable
          prop="website"
          label="គេហទំព័រ"
        />
        <AppInput
          v-model="form.facebook"
          placeholder="បញ្ចូលតំណភ្ជាប់ Facebook (មិនចាំបាច់)"
          clearable
          prop="facebook"
          label="ហ្វេសប៊ុក"
        />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.campuse-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>