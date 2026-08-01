// Utility to export array of rows to Excel-compatible CSV format
export const exportToCSV = (filename, rows) => {
  const processRow = (row) => {
    return row.map(val => {
      let result = val === null || val === undefined ? '' : String(val)
      result = result.replace(/"/g, '""')
      if (result.search(/("|,|\r|\n)/g) >= 0) {
        result = `"${result}"`
      }
      return result
    }).join(';')
  }

  const csvContent = '\uFEFF' + rows.map(processRow).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

export const formatDateCsv = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
