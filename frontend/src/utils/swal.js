import Swal from 'sweetalert2'

// Custom SweetAlert2 theme matching dark glassmorphism UI
const CustomSwal = Swal.mixin({
  background: '#1e293b',
  color: '#f8fafc',
  confirmButtonColor: '#6366f1',
  cancelButtonColor: '#94a3b8',
  customClass: {
    popup: 'swal-custom-popup',
    title: 'swal-custom-title',
    htmlContainer: 'swal-custom-text',
    confirmButton: 'swal-custom-confirm',
    cancelButton: 'swal-custom-cancel'
  }
})

export const alertSuccess = (title, text = '') => {
  return CustomSwal.fire({
    icon: 'success',
    title,
    text,
    timer: 2500,
    showConfirmButton: false
  })
}

export const alertError = (title, text = '') => {
  return CustomSwal.fire({
    icon: 'error',
    title,
    text,
    confirmButtonText: 'Tutup'
  })
}

export const alertWarning = (title, text = '') => {
  return CustomSwal.fire({
    icon: 'warning',
    title,
    text,
    confirmButtonText: 'Mengerti'
  })
}

export const confirmDialog = async (title, text, confirmButtonText = 'Ya, Lanjutkan') => {
  const result = await CustomSwal.fire({
    icon: 'question',
    title,
    text,
    showCancelButton: true,
    confirmButtonText,
    cancelButtonText: 'Batal',
    reverseButtons: true
  })
  return result.isConfirmed
}

export default CustomSwal
