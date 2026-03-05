import type { AssessmentSubmission, Assessment } from "../../types/assessment.d.ts";
// API Logic
export const useConsultation = () => {
  const config = useRuntimeConfig();
  const baseUri = config.public.apiBaseUri;

  const uploadAudio = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    const response = await $fetch<{
      status: string;
      message: string;
      statusCode: number;
      jobID?: string;
    }>(`${baseUri}/consultation/upload`, {
      method: "POST",
      body: formData,
    });
    return response;
  };

  const streamResult = async (
    jobId: string,
    onResult: (data: any) => void,
    onError: (error: any) => void,
  ) => {
    const eventSource = new EventSource(
      `${baseUri}/consultation/stream/${jobId}`,
    );

    eventSource.addEventListener("result", (e) => {
      if (!e.data) return
      const data = JSON.parse(e.data)
      onResult(data)
      eventSource.close()
    });

    eventSource.addEventListener("error", (e: MessageEvent) => {
      if (e.data) {
        try {
          const data = JSON.parse(e.data)
          onError(data)
        } catch {
          onError({ message: 'An error occurred during processing' })
        }
      } else {
        onError({ message: 'Could not connect to the processing stream. Please try again.' })
      }
      eventSource.close()
    });
    return eventSource;
  };

  const submitAssessment = async (submission: AssessmentSubmission) => {
    const response = await $fetch<{
      status: string;
      message: string;
      statusCode: number;
      data?: any;
    }>(`${baseUri}/consultation/collate-results`, {
      method: "POST",
      body: submission,
    })

    return response;
  }

  const getAssessment = async (id: string) => {
    const response = await $fetch<{
      status: string;
      message: string;
      statusCode: number;
      data?: any;
    }>(`${ baseUri }/assessments/${ id }`)

    return response;
  }

  const updateAssessment = async (id: string, submission: Assessment ) => {
    const response = await $fetch<{
      status: string;
      message: string;
      statusCode: number;
      data?: any;
    }>(`${ baseUri }/assessments/${ id }`, {
      method: "PUT",
      body: submission,
    })

    return response
  }
  return { uploadAudio, streamResult, submitAssessment, getAssessment, updateAssessment };
};
