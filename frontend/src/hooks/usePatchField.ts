export function usePatchField(reportID: string){
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (patch: PatchInput) => ApplyPatch(patch), // Wails Go binding
      onMount: async (patch) => {
      // Cancel any concurrent refetches for this reportID
      await queryClient.cancelQueries({queryKy: ['report', reportID]})
      // Snapshot current state
      const snapshot = queryClient.getQueryData<ReportData>(['report', reportID])
      // Optimistic update - user sees changes immediately
      queryClient.setQueryData(['report', reportID], (old: ReportData) => applyJsonPatch(old, patch))
     return {snapshot} 
    },
    onError: (_err, patch, context) => {
      // Server rejected -roll back to snapshot
      if (context?.snapshot) {
        queryClient.setQueryData(['report', reportID], context.snapshot)
      }
      showFieldError(patch.fieldPath, _err.message)
    },
  })
}
