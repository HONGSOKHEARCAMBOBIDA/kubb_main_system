  <script setup>
  import { ref, reactive, onMounted, watch } from 'vue'
  import {
    CreateClassCurriculum,
    GetClassCurriculum,
    UpdateClassCurriculumn,
    DeleteClassCurriculumn
  } from '../../services/classcurriculmn.service.js'
  import { getAcademics } from '../../services/academic.service.js'
  import { getGenerationByAcademic } from '../../services/generation.service.js'
  import { getTermByGeneation } from '../../services/term.service'
  import { getSemesterByAcademic } from '../../services/semester.service.js'
  import { getprogrammes } from '../../services/programmes.service.js'
  import { getFacultyByProgrammes } from '../../services/faculty.service.js'
  import { getDepartmentByFaculty } from '../../services/department.service.js'
  import { getMajorByDepartment, updateMajor } from '../../services/major.service.js'
  import { getAcademicshiftByAcademic } from '../../services/academic_shift.service.js'
  import { useNotification } from '../../composables/useNotification'
  import AppButton from '../../components/button/AppButton.vue'
  import AppInput from '../../components/input/AppInput.vue'
  import AppSelect from '../../components/common/AppSelect.vue'
  import AppForm from '@/components/forms/AppForm.vue'
  import AppDialog from '@/components/dialogs/AppDialog.vue'
  import AppFilterBar from '../../components/common/AppFilterBar.vue'
  import TableCustom from '../../components/tables/TableCustom.vue'
  import { CreateClassOffering } from '../../services/class_offering.service.js'
  import { getSubjectByMajor } from '../../services/subject.service.js'
  import { studentTermGet } from '../../services/studentterm.service.js'
  import { createcourseregistration } from '../../services/course_registration.service.js'
  import { getSubjectGroup } from '../../services/subject.service.js'
  const notify = useNotification()
  const classcurriculmn = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const submitting = ref(false)
  const formRef = ref(null)
  const dialogVisible = ref(false)
  const isEditMode = ref(false)
  const subjectgroup = ref([])



  // updateclasscurriculumn
  const updateuuid = ref(null)
  const updateacademicid = ref(null)
  const updategenerationid = ref(null)
  const updatetermid = ref(null)
  const updateprogrammid = ref(null)
  const updatefacultyid = ref(null)
  const updatedepartmentid = ref(null)
  const updatemajorid = ref(null)
  const formupdateclasscurricolumn = reactive({
    name: '',
    term_id: null,
    major_id: null
  })
  const updateclasscurriculumnvisible = ref(false)

  async function openupdateclasscurriculumn(row){
    updateuuid.value = row.uuid
    formupdateclasscurricolumn.name = row.name
    if(row.programme_id) {
      updateprogrammid.value = row.programme_id
      formFacultyOptions.value = await loadFacultyOption(row.programme_id)
    }
    if(row.faculty_id){
      updatefacultyid.value = row.faculty_id
      formDepartmentOptions.value = await loadDepartmentOption(row.faculty_id)
    }
    if(row.department_id) {
      updatedepartmentid.value = row.department_id
      formMajorOptions.value  = await loadMajorOption(row.department_id)
    }
    if (row.academic_id) {
      updateacademicid.value = row.academic_id
      generationOptions.value = await loadGenerationOptions(row.academic_id)
    }
    if (row.generation_id) {
      updategenerationid.value = row.generation_id
      termOptions.value = await loadTermOptions(row.generation_id)
    }
    formupdateclasscurricolumn.term_id = row.term_id ?? null
    formupdateclasscurricolumn.major_id = row.major_id ?? null
    updateclasscurriculumnvisible.value = true
  }

  async function updateacademicidchange(updateacademicid){
    updategenerationid.value = null
    formupdateclasscurricolumn.term_id = null
    generationOptions.value = []
    termOptions.value = []
    generationOptions.value = await loadGenerationOptions(updateacademicid)
  }

  async function updategenerationidchange(updategenerationid){
    formupdateclasscurricolumn.term_id = null
    termOptions.value = []
    termOptions.value = await loadTermOptions(updategenerationid)
  }

  async function updateprogrammidchange(updateprogrammid){
    updatefacultyid.value = null
    updatedepartmentid.value = null
    formupdateclasscurricolumn.major_id = null
    formFacultyOptions.value = []
    formDepartmentOptions.value = []
    formMajorOptions.value = []
    formFacultyOptions.value = await loadFacultyOption(updateprogrammid)
  }

async function updatefacultyidchange(updatefacultyid) {
     updatedepartmentid.value = null
    formupdateclasscurricolumn.major_id = null 
        formDepartmentOptions.value = []
    formMajorOptions.value = []
    formDepartmentOptions.value = await loadDepartmentOption(updatefacultyid)
}

async function updatedepartmentidchange(updatedepartmentid){
  formupdateclasscurricolumn.major_id = null 
   formMajorOptions.value = []
   formMajorOptions.value = await loadMajorOption(updatedepartmentid)
}

  async function handleSubmitUpdateClassCurriculum() {
    submitting.value = true
    try {
      await UpdateClassCurriculumn(updateuuid.value,formupdateclasscurricolumn)
      notify.success('កែប្រែកម្មវិធីសិក្សាបានជោគជ័យ')
      updateclasscurriculumnvisible.value = false
      resetForm()
      fetchClassCurriculum()
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to create')
    } finally {
      submitting.value = false
    }
  }

  // updateclasscurriculumn

  // deleteclasscurriculumn
  async function handleSubmitDeleteClassCurriculum(row) {
    submitting.value = true
    try {
      await DeleteClassCurriculumn(row.uuid)
      notify.success('លុបកម្មវិធីសិក្សាបានជោគជ័យ')
      fetchClassCurriculum()
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to create')
    } finally {
      submitting.value = false
    }
  }
  // deleteclasscurriculumn

  const studentcolumns = [
  { prop: 'name_kh',slot:'name_kh', label: 'ឈ្មោះ', minwidth: 120 }, 
  { prop: 'date_of_birth', label: 'ថ្ងៃ-ខែ-ឆ្នាំកំណើត', minwidth: 120 }, 
  { prop: 'gender',slot:'gender', label: 'ភេទ', minwidth: 120 }, 
  { prop: 'nationality', label: 'សញ្ជាតិ', minwidth: 120 },
  { prop: 'phone', label: 'លេខទូរសព្ទ', minwidth: 120 },
  { prop: 'occupation', label: 'មុខរបរ', minwidth: 120 },
  ]

  const columns = [
    { prop: 'name', label: 'ឈ្មោះថ្នាក់', minwidth: 120 },
    { prop: 'major_name', slot: 'major_name', label: 'ជំនាញ', minwidth: 120 },
        { prop: 'generation_name', label: 'ជំនាន់', width: 120 },
    { prop: 'term_name', label: 'វគ្គ', width: 120 },
      { slot: 'programme_name', label: 'កម្រិត', width: 140 },
    {
      prop: 'major_duration_interval',
      slot: 'major_duration_interval',
      label: 'រយៈពេលសិក្សា',
      width: 120,
    },
    { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', width: 120 },
    { prop: 'active', slot: 'active', label: 'ស្ថានភាព', width: 120 },
  ]

  const columndetails = [
    { prop: 'semester_name', label: 'ឆមាស', minwidth: 120 },
    { prop: 'study_year_id', label: 'ឆ្នាំទី', minwidth: 120 },
    { prop: 'academic_name', label: 'ឆ្នាំសិក្សា', minwidth: 120 },
    { prop: 'academic_shift_name', label: 'វេនសិក្សា', minwidth: 120 },
    { prop: 'midterm_date', label: 'ថ្ងៃប្រឡង Midterm', minwidth: 120 },
    { prop: 'final_date', label: 'ថ្ងៃប្រឡង Final', minwidth: 120 },
    { prop: 'total_student', label: 'ប្រធានថ្នាក់', minwidth: 120 },
    { prop: 'type_class', slot: 'type_class', label: 'ប្រភេទថ្នាក់', minwidth: 120 },
  ]

  const columnclass_offering = [
    { prop: 'subject_code', label: 'កូដ', minwidth: 90 },
    { prop: 'subject_name',slot:'subject_name', label: 'ឈ្មោះ', minminwidth: 120 },
    {prop:'subject_group_name',slot:'subject_group_name',label:'ក្រុម',width: 120},
    { prop: 'credit',slot:'credit', label: 'ក្រេឌីត', minwidth: 90 },
    { prop: 'passing_score',slot:'passing_score', label: 'ពិន្ទុជាប់', minwidth: 90 },
    { prop: 'total_hour',slot:'total_hour', label: 'ម៉ោងសរុប', minwidth: 90 },
    { prop: 'total_hour',slot:'total_hour', label: 'ម៉ោងនៅសល់', minwidth: 140 },
    { prop: 'status', label: 'ស្ថានភាព', minwidth: 90 },
    { prop: 'total_attendance_for_rexam',slot:'total_attendance_for_rexam', label: 'អវត្តមានប្រឡងសង', minwidth: 120 },
    { prop: 'total_attendance_for_relearn',slot:'total_attendance_for_relearn', label: 'អវត្តមានរៀនសង', minwidth: 120 },
    { prop: 'description', label: 'ផ្សេងៗ', minwidth: 170 },
  ]

  async function fetchClassCurriculum() {
    try {
      const params = {
        page: page.value,
        page_size: pageSize.value,
      }
      const res = await GetClassCurriculum(params)
      classcurriculmn.value = res.data.data || []
      total.value = res.data.pagination.total_count || 0
      console.log(classcurriculmn.value)
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load class curriculums')
    }
  }

  /* ---------------- Class Offering dialog ---------------- */
  const createclassofferingvisible = ref(false)
  const submittingOffering = ref(false)
  const subjectoptions = ref([])

  function emptyOfferingRow() {
    return {
      subject_id: null,
      subject_group_id:null,
      credit: 0,
      passing_score: 0,
      total_hour: 0,
      total_attendance_for_rexam: 0,
      total_attendance_for_relearn: 0,
      description: '',
    }
  }

  function defaultOfferingForm() {
    return {
      class_curriculum_detail_id: null,
      class_offering: [emptyOfferingRow()],
    }
  }

  const classofferingform = reactive(defaultOfferingForm())

  function newClassOffering() {
    classofferingform.class_offering.push(emptyOfferingRow())
  }

  function removeClassOfferingRow(index) {
    if (classofferingform.class_offering.length <= 1) return
    classofferingform.class_offering.splice(index, 1)
  }

  // detailRow = the row from class_curriculum_detais (has the detail id)
  // curriculumRow = the parent class curriculum row (has major_id)
  async function openClassOfferingDialog(detailRow, curriculumRow) {
    Object.assign(classofferingform, defaultOfferingForm())
    classofferingform.class_curriculum_detail_id = detailRow.id
    subjectoptions.value = await loadSubjectOption(curriculumRow.major_id)
    createclassofferingvisible.value = true
  }

  function closeClassOfferingDialog() {
    createclassofferingvisible.value = false
    subjectoptions.value = []
    Object.assign(classofferingform, defaultOfferingForm())
  }

  async function submitClassOffering() {
    submittingOffering.value = true
    try {
      await CreateClassOffering(classofferingform)
      notify.success('បញ្ចូលមុខវិជ្ជាបានជោគជ័យ')
      closeClassOfferingDialog()
      fetchClassCurriculum()
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to create class offering')
    } finally {
      submittingOffering.value = false
    }
  }

  /* ---------------- Class Curriculum top-level form ---------------- */


  // class registration start
  const createclassregistrationvisible = ref(false)
  const selectedRows = ref([]);
  const studentterms = ref([])
  const selectedClassOffering = ref(null)

  const classregistrationfilters = reactive({
      semester_id: null,
      study_year_id: null,
      term_id: null,
      major_id: null
  })

  const studentgradecolumn = [
    { prop: 'total_score',label: 'Total Score', minwidth: 120 },
    { prop: 'letter_grade',label: 'Letter Grade', minwidth: 120 },
    { prop: 'grade_point',label: 'Grade Point', minwidth: 120 },
  ]

  const studentgradedetailcolumn = [
    { prop: 'grade_component_name',label: 'ប្រភេទ', minwidth: 120 },
    { prop: 'score',label: 'ពិន្ទុ', minwidth: 120 },
  ]

  const columnclass_registration = [
    { prop: 'student_name_kh',slot:'student_name_kh', label: 'ឈ្មោះ', minwidth: 90 },
    { prop: 'student_gender',slot:'student_gender',label: 'ភេទ', minwidth: 90 },
    { prop: 'study_year_id',label: 'ឆ្នាំទី', minwidth: 90 },
    { prop: 'semester_name',label: 'ឆមាសទី', minwidth: 90 },
    { prop: 'term_name',label: 'វគ្គទី', minwidth: 90 },
    { prop: 'major_name',slot:'major_name',label: 'ជំនាញ', minwidth: 90 },
      { prop: 'programm_name',slot:'programm_name',label: 'កម្រិត', minwidth: 90 },
  ]

  async function fetchStudentTerm(filters) {
      try {
          const params = {}

          if (filters.semester_id) {
              params.semester_id = filters.semester_id
          }

          if (filters.study_year_id) {
              params.study_year_id = filters.study_year_id
          }

          if (filters.term_id) {
              params.term_id = filters.term_id
          }

          if (filters.major_id) {
              params.major_id = filters.major_id
          }

          const res = await studentTermGet(params)

          return res.data.data || []
      } catch (e) {
          notify.error(
              e?.response?.data?.message ||
              e.message ||
              'Failed to load students'
          )

          return []
      }
  }

  async function openClassRegistrationDialog(
    detailRow,
    curriculumRow,
    class_offering
  ) {
    classregistrationfilters.semester_id = detailRow.semester_id
    classregistrationfilters.study_year_id = detailRow.study_year_id
    classregistrationfilters.term_id = curriculumRow.term_id
    classregistrationfilters.major_id = curriculumRow.major_id

    selectedClassOffering.value = class_offering
    selectedRows.value = []

    studentterms.value = await fetchStudentTerm(
      classregistrationfilters
    )

    createclassregistrationvisible.value = true
  }

  function closeClassRegistrationDialog() {
    createclassregistrationvisible.value = false
    studentterms.value = []
  }

  async function submitcourseregistration() {
    if (!selectedClassOffering.value?.id) {
      notify.error('Class offering is required')
      return
    }

    if (selectedRows.value.length === 0) {
      notify.error('Please select at least one student')
      return
    }

    try {
      const payload = {
        class_offering_id: selectedClassOffering.value.id,
        student_term_id: selectedRows.value.map((row) => ({
          student_term_id: row.id,
        })),
      }

      await createcourseregistration(payload)

      notify.success('ចុះឈ្មោះសិស្សបានជោគជ័យ')

      selectedRows.value = []
      closeClassRegistrationDialog()
    } catch (e) {
      notify.error(
        e?.response?.data?.message ||
        e.message ||
        'Failed to create course registration'
      )
    }
  }

  // class registration end
  function emptyDetailRow() {
    return {
      semester_id: null,
      study_year_id: null,
      academic_shift_id: null,
      midterm_date: '',
      final_date: '',
      type_class: '',
    }
  }

  function defaultForm() {
    return {
      name: '',
      major_id: null,
      term_id: null,
      class_curriclumn_details: [emptyDetailRow()],
    }
  }

  const form = reactive(defaultForm())

  function addDetailRow() {
    form.class_curriclumn_details.push(emptyDetailRow())
  }

  function removeDetailRow(index) {
    if (form.class_curriclumn_details.length <= 1) return
    form.class_curriclumn_details.splice(index, 1)
  }

  const rules = {
    name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះកម្មវិធីសិក្សា', trigger: 'blur' }],
    major_id: [{ required: true, message: 'សូមជ្រើសរើសជំនាញ', trigger: 'change' }],
    term_id: [{ required: true, message: 'សូមជ្រើសរើសវគ្គ', trigger: 'change' }],
  }

  /* cascading Program -> Faculty -> Department -> Major */
  const formProgramID = ref(null)
  const formFacultyID = ref(null)
  const formDepartmentID = ref(null)
  const programmesOptions = ref([])
  const formFacultyOptions = ref([])
  const formDepartmentOptions = ref([])
  const formMajorOptions = ref([])

  /* cascading Academic -> Generation -> Term, and Academic -> Semester / Academic Shift */
  const formAcademicId = ref(null)
  const formGenerationId = ref(null)
  const academicOptions = ref([])
  const generationOptions = ref([])
  const termOptions = ref([])
  const semesterOptions = ref([])
  const academicShiftOptions = ref([])

  const studyyearOption = [
    { label: 'ឆ្នាំទី1', value: 1 },
    { label: 'ឆ្នាំទី2', value: 2 },
    { label: 'ឆ្នាំទី3', value: 3 },
    { label: 'ឆ្នាំទី4', value: 4 },
  ]

  const typeClassOptions = [
    { label: 'ថ្នាក់ផ្ទាល់', value: 'onclass' },
    { label: 'ថ្នាក់អនឡាញ', value: 'online' },
  ]

  /* ---------------- option loaders ---------------- */

  async function fetchSubjectGroupOptions() {
    try {
      const res = await getSubjectGroup()
      subjectgroup.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load programmes')
    }
  }

  async function fetchProgrammeOptions() {
    try {
      const res = await getprogrammes()
      programmesOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load programmes')
    }
  }

  async function loadFacultyOption(programmeID) {
    if (!programmeID) return []
    try {
      const res = await getFacultyByProgrammes(programmeID)
      return (res.data.data || []).map((f) => ({ label: f.name, value: f.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load faculties')
      return []
    }
  }

  async function loadDepartmentOption(facultyID) {
    if (!facultyID) return []
    try {
      const res = await getDepartmentByFaculty(facultyID)
      return (res.data.data || []).map((d) => ({ label: d.name, value: d.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load departments')
      return []
    }
  }

  async function loadMajorOption(departmentID) {
    if (!departmentID) return []
    try {
      const res = await getMajorByDepartment(departmentID)
      return (res.data.data || []).map((m) => ({ label: m.name, value: m.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load majors')
      return []
    }
  }

  async function loadSubjectOption(majorID) {
    if (!majorID) return []
    try {
      const res = await getSubjectByMajor(majorID)
      return (res.data.data || []).map((m) => ({ label: m.name, value: m.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load subjects')
      return []
    }
  }

  async function fetchAcademicOptions() {
    try {
      const res = await getAcademics()
      academicOptions.value = (res.data.data || []).map((a) => ({ label: a.name, value: a.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load academics')
    }
  }

  async function loadGenerationOptions(academicId) {
    if (!academicId) return []
    try {
      const res = await getGenerationByAcademic(academicId)
      return (res.data.data || []).map((g) => ({ label: g.name, value: g.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load generations')
      return []
    }
  }

  async function loadTermOptions(generationId) {
    if (!generationId) return []
    try {
      const res = await getTermByGeneation(generationId)
      return (res.data.data || []).map((t) => ({ label: t.name, value: t.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load terms')
      return []
    }
  }

  async function loadSemesterOptions(academicId) {
    if (!academicId) return []
    try {
      const res = await getSemesterByAcademic(academicId)
      return (res.data.data || []).map((s) => ({ label: s.name, value: s.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load semesters')
      return []
    }
  }

  async function loadAcademicShiftOptions(academicId) {
    if (!academicId) return []
    try {
      const res = await getAcademicshiftByAcademic(academicId)
      return (res.data.data || []).map((s) => ({ label: s.name, value: s.id }))
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to load academic shifts')
      return []
    }
  }

  /* Program -> Faculty -> Department -> Major */
  watch(formProgramID, async (newVal) => {
    formFacultyID.value = null
    formDepartmentID.value = null
    formDepartmentOptions.value = []
    formMajorOptions.value = []
    form.major_id = null
    formFacultyOptions.value = await loadFacultyOption(newVal)
  })

  watch(formFacultyID, async (newVal) => {
    formDepartmentID.value = null
    formMajorOptions.value = []
    form.major_id = null
    formDepartmentOptions.value = await loadDepartmentOption(newVal)
  })

  watch(formDepartmentID, async (newVal) => {
    form.major_id = null
    formMajorOptions.value = await loadMajorOption(newVal)
  })

  /* Academic -> Generation, Semester, Academic Shift */
  watch(formAcademicId, async (newVal) => {
    formGenerationId.value = null
    form.term_id = null
    generationOptions.value = []
    termOptions.value = []
    semesterOptions.value = []
    academicShiftOptions.value = []

    if (!newVal) return

    generationOptions.value = await loadGenerationOptions(newVal)
    semesterOptions.value = await loadSemesterOptions(newVal)
    academicShiftOptions.value = await loadAcademicShiftOptions(newVal)
  })

  /* Generation -> Term */
  watch(formGenerationId, async (newVal) => {
    form.term_id = null
    termOptions.value = []

    if (!newVal) return

    termOptions.value = await loadTermOptions(newVal)
  })

  /* ---------------- dialog open / close / submit ---------------- */

  function openCreate() {
    isEditMode.value = false
    resetForm()
    dialogVisible.value = true
  }

  function closeDialog() {
    dialogVisible.value = false
    resetForm()
  }

  async function handleSubmit() {
    submitting.value = true
    try {
      await CreateClassCurriculum(form)
      notify.success('បង្កើតកម្មវិធីសិក្សាបានជោគជ័យ')
      dialogVisible.value = false
      resetForm()
      fetchClassCurriculum()
    } catch (e) {
      notify.error(e?.response?.data?.message || e.message || 'Failed to create')
    } finally {
      submitting.value = false
    }
  }



  function resetForm() {
    Object.assign(form, defaultForm())

    formProgramID.value = null
    formFacultyID.value = null
    formDepartmentID.value = null
    formFacultyOptions.value = []
    formDepartmentOptions.value = []
    formMajorOptions.value = []

    formAcademicId.value = null
    formGenerationId.value = null
    generationOptions.value = []
    termOptions.value = []
    semesterOptions.value = []
    academicShiftOptions.value = []

    formRef.value?.resetFields?.()
  }

  onMounted(() => {
    fetchSubjectGroupOptions()
    fetchClassCurriculum()
    fetchProgrammeOptions()
    fetchAcademicOptions()
  })
  </script>

  <template>
    <AppFilterBar :fields="[{ slot: 'create', span: 4 }]">
      <template #create>
        <AppButton type="default" icon="Plus" @click="openCreate">បង្កើតថ្នាក់ថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      expandable
      :data="classcurriculmn"
      :columns="columns"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchClassCurriculum"
    >
      <template #programme_name="{ row }">
        <el-text>
          {{ row.programme_name }}
        </el-text>
      </template>
      <template #active="{ row }">
        <el-text>
          {{ row.active === true ? 'កំពុងសិក្សា' : 'បានបញ្ចប់' }}
        </el-text>
      </template>
      <template #major_name="{ row }">
        <el-text > {{ row.major_code }} | <el-text size="small" type="primary">{{ row.major_name }}</el-text> | <el-text size="small">{{ row.faculty_name }}</el-text> </el-text>
      </template>
      <template #major_duration_interval="{ row }">
        <el-text > {{ row.major_duration_period }} {{ row.major_duration_interval }} </el-text>
      </template>
      <template #actions="{row}">
        <el-tooltip content="កែប្រែ" placement="top">
        <AppButton
        circle
        plain
        type="warning"
        icon="Edit"
        size="small"
        @click="openupdateclasscurriculumn(row)"
        >

        </AppButton>
        </el-tooltip>
        <el-tooltip content="លុប" placement="top">
        <AppButton
        circle
        plain
        type="danger"
        icon="Delete"
        size="small"
        @click="handleSubmitDeleteClassCurriculum(row)"
        >

        </AppButton>
        </el-tooltip>
      </template>
      <!-- `row` here is the parent class curriculum row -->
      <template #expand="{ row }">
        <el-divider content-position="left">
          <AppButton
          plain
          >
            ថែម
          </AppButton>
        </el-divider>
        <TableCustom
          expandable
          :data="row.class_curriculum_detais"
          :columns="columndetails"
          :show-pagination="false"
        >
          <template #type_class="{ row }">
            <el-text>
              {{ row.type_class === 'onclass' ? 'រៀនថ្នាក់ផ្ទាល់' : 'រៀនOnline' }}
            </el-text>
          </template>
          <!-- rename inner scope's `row` to `detailRow` so it doesn't shadow the outer curriculum `row` -->
          <template #actions="{ row: detailRow }">
            <AppButton @click="openClassOfferingDialog(detailRow, row)"> ថែមមុខវិជ្ជា </AppButton>
          </template>
          <template #expand="{  row: detailRow }">
            <el-divider content-position="left">
            <el-text tag="b" style="color: darkcyan;">
              មុខវិជ្ជាត្រូវសិក្សា
            </el-text>
            </el-divider>
            <TableCustom
              expandable
              :data="detailRow.class_offering"
              :columns="columnclass_offering"
              :show-pagination="false"
              actions-width="200px"
            >
            <template #credit="{row}">
              <el-text >
                {{ row.credit }}
              </el-text>
            </template>
            <template #subject_group_name="{row}">
              <el-text size="small" type="primary">
                {{ row.subject_group_name }}
              </el-text>
            </template>
            <template #passing_score="{row}">
              <el-text>
                {{ row.passing_score }}
              </el-text>
            </template>
            <template #total_hour="{row}">
              <el-text >
                {{ row.total_hour }} ម៉ោង
              </el-text>
            </template>
            <template #subject_name="{row}">
              <el-text>
                {{ row.subject_name }}
              </el-text>
            </template>
            <template #total_attendance_for_rexam="{row}">
              <el-text >
                {{ row.total_attendance_for_rexam }}
              </el-text>
            </template>
            <template #total_attendance_for_relearn="{row}">
              <el-text >
                {{ row.total_attendance_for_relearn }}
              </el-text>
            </template>
      <template #actions="{row: class_offering}">
        <AppButton
        size="small"
          type="success"
          plain
          @click="openClassRegistrationDialog(detailRow, row, class_offering)"
        >
          បញ្ចូលសិស្ស
        </AppButton>
        <AppButton
        size="small"
          type="primary"
          plain
        >
          កែប្រែ
        </AppButton>
      </template>
  <template #expand="{ row: offeringRow }">
    <el-divider content-position="left">
      <el-text tag="b" style="color: darkcyan;">សិស្សសរុបកំពុងសិក្សា {{ offeringRow.student.length }} នាក់</el-text>
    </el-divider>
    <TableCustom
      expandable
      :data="offeringRow.student"
      :columns="studentcolumns"
      :show-pagination="false"
    >
    <template #name_kh="{row}">
      <el-text>
        {{ row.name_kh }} | <el-text type="primary" size="small">{{ row.name_en }}</el-text>
      </el-text>
    </template>
    <template #gender="{row}">
      <el-text>
        {{ row.gender === 'Male' ? 'ប្រុស' : 'ស្រី' }}
      </el-text>
    </template>

    <template #expand="{ row: studentgrade }">
      <el-divider>
        Student Grade
      </el-divider>
      <TableCustom
      expandable
      :data="studentgrade.student_grade"
      :columns="studentgradecolumn"
      :show-pagination="false"
      >

      <template #expand="{ row: studentgradedetail }">
      <el-divider>
        Detail
      </el-divider>
      <TableCustom
      expandable
      :data="studentgradedetail.detail"
      :columns="studentgradedetailcolumn"
      :show-pagination="false"
      >

      </TableCustom>
      </template>

      </TableCustom>
    </template>

  </TableCustom>
  </template>
      
            </TableCustom>
          </template>
        
        </TableCustom>
        
      </template>
    </TableCustom>

    <AppDialog
      v-if="dialogVisible"
      v-model:visible="dialogVisible"
      title="បង្កើតកម្មវិធីសិក្សាថ្មី"
      :showDefaultFooter="false"
      width="60%"
      @close="closeDialog"
    >
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :loading="submitting"
        :show-actions="true"
        @submit="handleSubmit"
        submitText="រក្សាទុក"
        resetText="ចាកចេញ"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model="form.name"
              placeholder="បញ្ចូលឈ្មោះកម្មវិធីសិក្សា"
              clearable
              prop="name"
              label="ឈ្មោះកម្មវិធីសិក្សា"
            />
          </el-col>
        </el-row>

        <!-- cascading select: Academic -> Generation -> Term -->
        <el-row :gutter="20">
          <el-col :span="8">
            <AppSelect
              v-model="formAcademicId"
              :options="academicOptions"
              placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
              clearable
              label="ឆ្នាំសិក្សា"
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="formGenerationId"
              :options="generationOptions"
              placeholder="ជ្រើសរើសជំនាន់"
              clearable
              label="ជំនាន់"
              :disabled="!formAcademicId"
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="form.term_id"
              :options="termOptions"
              placeholder="ជ្រើសរើសវគ្គ"
              clearable
              prop="term_id"
              label="វគ្គ"
              :disabled="!formGenerationId"
            />
          </el-col>
        </el-row>

        <!-- cascading select down to major -->
        <el-row :gutter="20">
          <el-col :span="6">
            <AppSelect
              v-model="formProgramID"
              :options="programmesOptions"
              placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
              clearable
              label="កម្មវិធីសិក្សា"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="formFacultyID"
              :options="formFacultyOptions"
              placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
              clearable
              label="មហាវិទ្យាល័យ"
              :disabled="!formProgramID"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="formDepartmentID"
              :options="formDepartmentOptions"
              placeholder="ជ្រើសរើសដេប៉ាតឺម៉ង់"
              clearable
              label="ដេប៉ាតឺម៉ង់"
              :disabled="!formFacultyID"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="form.major_id"
              :options="formMajorOptions"
              placeholder="ជ្រើសរើសជំនាញ"
              clearable
              prop="major_id"
              label="ជំនាញ"
              :disabled="!formDepartmentID"
            />
          </el-col>
        </el-row>

        <!-- detail rows -->
        <div class="detail-section">
          <div class="detail-header">
            <el-text tag="b">ព័ត៌មានលម្អិត (ឆមាស / ឆ្នាំសិក្សា / វេន)</el-text>
            <AppButton type="default" icon="Plus" size="small" @click="addDetailRow"
              >ថែមជួរ</AppButton
            >
          </div>

          <div v-for="(row, index) in form.class_curriclumn_details" :key="index" class="detail-row">
            <div class="detail-row-header">
              <el-text tag="b" type="info">ឆមាសទី {{ index + 1 }}</el-text>
              <AppButton
                icon="Delete"
                circle
                size="small"
                type="default"
                plain
                :disabled="form.class_curriclumn_details.length <= 1"
                @click="removeDetailRow(index)"
              />
            </div>

            <el-row :gutter="20">
              <el-col :span="8">
                <AppSelect
                  v-model="row.study_year_id"
                  :options="studyyearOption"
                  placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
                  clearable
                  label="ឆ្នាំសិក្សា"
                />
              </el-col>
              <el-col :span="8">
                <AppSelect
                  v-model="row.semester_id"
                  :options="semesterOptions"
                  placeholder="ជ្រើសរើសឆមាស"
                  clearable
                  label="ឆមាស"
                  :disabled="!formAcademicId"
                />
              </el-col>
              <el-col :span="8">
                <AppSelect
                  v-model="row.academic_shift_id"
                  :options="academicShiftOptions"
                  placeholder="ជ្រើសរើសវេន"
                  clearable
                  label="វេន"
                  :disabled="!formAcademicId"
                />
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="8">
                <AppInput
                  v-model="row.midterm_date"
                  type="date"
                  placeholder="ថ្ងៃប្រឡង Midterm"
                  clearable
                  label="ថ្ងៃប្រឡង Midterm"
                />
              </el-col>
              <el-col :span="8">
                <AppInput
                  v-model="row.final_date"
                  type="date"
                  placeholder="ថ្ងៃប្រឡង Final"
                  clearable
                  label="ថ្ងៃប្រឡង Final"
                />
              </el-col>
              <el-col :span="8">
                <AppSelect
                  v-model="row.type_class"
                  :options="typeClassOptions"
                  placeholder="ជ្រើសរើសប្រភេទថ្នាក់"
                  clearable
                  label="ប្រភេទថ្នាក់"
                />
              </el-col>
            </el-row>
          </div>
        </div>
      </AppForm>
    </AppDialog>

    <!-- Class Offering dialog -->
    <AppDialog
      v-if="createclassofferingvisible"
      v-model:visible="createclassofferingvisible"
      title="បញ្ចូលមុខវិជ្ជាទៅថ្នាក់រៀន"
      :showDefaultFooter="false"
      width="720px"
      @close="closeClassOfferingDialog"
    >
      <AppForm
        :model="classofferingform"
        :loading="submittingOffering"
        :show-actions="true"
        @submit="submitClassOffering"
        submitText="បញ្ជូល"
      >
        <div class="section-header">
          <el-text tag="b">បញ្ចូលមុខវិជ្ជា</el-text>
          <AppButton type="warning" plain icon="Plus" size="small" @click="newClassOffering">
            បន្ថែម
          </AppButton>
        </div>

        <el-card
          v-for="(classf, index) in classofferingform.class_offering"
          :key="index"
          :gutter="16"
          style="margin-bottom: 12px"
        >
          <template #header>
            <div class="section-header">
              <span>មុខវិជ្ជាទី -{{ index + 1 }}</span>
              <AppButton
                type="danger"
                icon="Delete"
                size="small"
                plain
                :disabled="classofferingform.class_offering.length <= 1"
                @click="removeClassOfferingRow(index)"
              >
                លុប
              </AppButton>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :span="8">
              <AppSelect
                v-model="classf.subject_id"
                :options="subjectoptions"
                placeholder="មុខវិជ្ជា"
                label="មុខវិជ្ជា"
                clearable
              />
            </el-col>
            <el-col :span="8">
              <AppSelect
                v-model="classf.subject_group_id"
                :options="subjectgroup"
                placeholder="ប្រភេទ"
                label="ប្រភេទ"
                clearable
              />
            </el-col>
            <el-col :span="8">
              <AppInput
                v-model.number="classf.credit"
                placeholder="បញ្ចូលក្រេឌីត"
                type="number"
                clearable
                label="បញ្ចូលក្រេឌីត"
              />
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <AppInput
                v-model.number="classf.passing_score"
                placeholder="បញ្ចូលពន្ទុជាប់"
                type="number"
                clearable
                label="បញ្ចូលពន្ទុជាប់"
              />
            </el-col>
            <el-col :span="12">
              <AppInput
                v-model.number="classf.total_hour"
                placeholder="បញ្ចូលចំនួនម៉ោងត្រូវរៀន"
                type="number"
                clearable
                label="បញ្ចូលចំនួនម៉ោងត្រូវរៀន"
              />
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <AppInput
                v-model.number="classf.total_attendance_for_relearn"
                placeholder="បញ្ចូលចំនួនអវត្តមានត្រូវរៀនសង"
                type="number"
                clearable
                label="បញ្ចូលចំនួនអវត្តមានត្រូវរៀនសង"
              />
            </el-col>
            <el-col :span="12">
              <AppInput
                v-model.number="classf.total_attendance_for_rexam"
                placeholder="បញ្ចូលចំនួនវត្តមានត្រូវប្រឡងសង"
                type="number"
                clearable
                label="បញ្ចូលចំនួនវត្តមានត្រូវប្រឡងសង"
              />
            </el-col>
          </el-row>

          <AppInput v-model="classf.description" placeholder="ផ្សេងៗ" clearable label="ផ្សេងៗ" />
        </el-card>
      </AppForm>
    </AppDialog>
    <AppDialog
      v-if="createclassregistrationvisible"
      v-model:visible="createclassregistrationvisible"
      title="បញ្ចូលសិស្ស"
      :showDefaultFooter="false"
      width="70%"
      @close="closeClassRegistrationDialog"
    >
      <AppForm
        :show-actions="true"
        @submit="submitcourseregistration"
        submitText="រក្សាទុក"
        resetText="ចាកចេញ"   
      >
          <TableCustom

    selectable
          :data="studentterms"
          :columns="columnclass_registration"
          :show-pagination="false"
    @selection-change="selectedRows = $event"
    >
    <template #student_name_kh="{row}">
      <div>
        <el-text tag="b" style="color: black;">
          {{ row.student_name_kh }}
        </el-text>
      </div>
      <div>
        <el-text type="primary">
          {{ row.student_name_en }}
        </el-text>
      </div>
    </template>
    <template #student_gender="{row}">
      <el-text style="color: black;">
        {{ row.student_gender === 'Male' ? 'ប្រុស' : 'ស្រី' }}
      </el-text>
    </template>
    <template #major_name="{row}">
      <el-text type="warning" tag="b">
      ({{ row.major_code }}) 
      </el-text>
      <el-text tag="b" style="color: darkcyan;">
        {{ row.major_name}}
      </el-text>
    </template>

    <template #programm_name="{row}">
      <el-text tag="b" style="color: crimson;">
        {{ row.programm_name }}
      </el-text>
    </template>

    </TableCustom>
      </AppForm>
    </AppDialog>

    <AppDialog
      v-if="updateclasscurriculumnvisible"
      v-model:visible="updateclasscurriculumnvisible"
      title="កែប្រែកម្មវិធីសិក្សាថ្មី"
      :showDefaultFooter="false"
      width="60%"
      @close="closeDialog"
    >
      <AppForm
        ref="formRef"
        :model="formupdateclasscurricolumn"
        :show-actions="true"
        @submit="handleSubmitUpdateClassCurriculum"
        submitText="រក្សាទុក"
        resetText="ចាកចេញ"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <AppInput
              v-model="formupdateclasscurricolumn.name"
              placeholder="បញ្ចូលឈ្មោះកម្មវិធីសិក្សា"
              clearable
              label="ឈ្មោះកម្មវិធីសិក្សា"
            />
          </el-col>
        </el-row>

        <!-- cascading select: Academic -> Generation -> Term -->
        <el-row :gutter="20">
          <el-col :span="8">
            <AppSelect
              v-model="updateacademicid"
              :options="academicOptions"
              placeholder="ជ្រើសរើសឆ្នាំសិក្សា"
              clearable
              label="ឆ្នាំសិក្សា"
              @change="updateacademicidchange"
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="updategenerationid"
              :options="generationOptions"
              placeholder="ជ្រើសរើសជំនាន់"
              clearable
              label="ជំនាន់"
              :disabled="!updateacademicid"
              @change="updategenerationidchange"
            />
          </el-col>
          <el-col :span="8">
            <AppSelect
              v-model="formupdateclasscurricolumn.term_id"
              :options="termOptions"
              placeholder="ជ្រើសរើសវគ្គ"
              clearable
              label="វគ្គ"
              :disabled="!updategenerationid"
            />
          </el-col>
        </el-row>

        <!-- cascading select down to major -->
        <el-row :gutter="20">
          <el-col :span="6">
            <AppSelect
              v-model="updateprogrammid"
              :options="programmesOptions"
              placeholder="ជ្រើសរើសកម្មវិធីសិក្សា"
              clearable
              label="កម្មវិធីសិក្សា"
              @change="updateprogrammidchange"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="updatefacultyid"
              :options="formFacultyOptions"
              placeholder="ជ្រើសរើសមហាវិទ្យាល័យ"
              clearable
              label="មហាវិទ្យាល័យ"
              :disabled="!updateprogrammid"
              @change="updatefacultyidchange"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="updatedepartmentid"
              :options="formDepartmentOptions"
              placeholder="ជ្រើសរើសដេប៉ាតឺម៉ង់"
              clearable
              label="ដេប៉ាតឺម៉ង់"
              :disabled="!updatefacultyid"
              @change="updatedepartmentidchange"
            />
          </el-col>
          <el-col :span="6">
            <AppSelect
              v-model="formupdateclasscurricolumn.major_id"
              :options="formMajorOptions"
              placeholder="ជ្រើសរើសជំនាញ"
              clearable
              label="ជំនាញ"
              :disabled="!updatedepartmentid"
            />
          </el-col>
        </el-row>

      </AppForm>
    </AppDialog>

  </template>

  <style scoped>
  .detail-section {
    margin-top: 16px;
    border-top: 1px solid #ebeef5;
    padding-top: 16px;
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .detail-row {
    border: 1px solid #ebeef5;
    border-radius: 4px;
    padding: 12px;
    margin-bottom: 12px;
    background: #fafafa;
  }

  .detail-row-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }
  </style>
