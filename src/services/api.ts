/**
 * API Service Layer
 * Handles all communication with the backend API
 */

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

/**
 * Base fetch wrapper with error handling
 */
async function apiFetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`API Error (${response.status}): ${errorText}`);
    }

    return response.json();
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    throw error;
  }
}
/**
 * Base fetch wrapper for requests that do not expect a JSON body
 * - good for DELETE, POST/PUT that only return 204, etc.
 */
async function apiFetchNoJSON(
  endpoint: string,
  options?: RequestInit
): Promise<void> {
  const url = `${API_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`API Error (${response.status}): ${errorText}`);
    }

    // Nothing to return - just resolve
    return;
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    throw error;
  }
}

// ==================
// IMAGE API
// ==================

export interface Image {
  id: string;
  slug?: string;
  object_url: string;
  name: string;
  text: string;
  date?: string;
  location?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateImageRequest {
  slug?: string;
  name: string;
  text: string;
  date?: string;
  location?: string;
  data: string; // Base64 encoded image
}

export interface UpdateImageRequest {
  name?: string;
  text?: string;
  date?: string;
  location?: string;
  data?: string; // Base64 encoded image (optional)
}

/**
 * Get image by ID
 */
export async function getImageById(id: string): Promise<Image> {
  return apiFetch<Image>(`/images/${id}`);
}

/**
 * List all images
 */
export async function listAllImages(): Promise<Image[]> {
  return apiFetch<Image[]>('/images');
}

/**
 * Get images by gallery slug
 */
// After
export async function getImagesBySlug(slug: string): Promise<Image[]> {
  const res = await apiFetch<Image[] | null>(`/images/slug/${slug}`);
  return res ?? [];
}


/**
 * Upload new image
 */
export async function createImage(request: CreateImageRequest): Promise<Image> {
  return apiFetch<Image>('/images', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}

/**
 * Update image metadata and/or data
 */
export async function updateImage(id: string, request: UpdateImageRequest): Promise<Image> {
  return apiFetch<Image>(`/images/${id}`, {
    method: 'PUT',
    body: JSON.stringify(request),
  });
}

/**
 * Delete image
 */
export async function deleteImage(id: string): Promise<void> {
  await apiFetchNoJSON(`/images/${id}`, {
    method: 'DELETE',
  });
}

// ==================
// GALERY EVENTS API
// ==================

export interface GaleryEvent {
  id: string;
  name: string;
  location: string;
  date: string; // ISO 8601 date string
  image_urls: string[];
  image_ids: string[];
  created_at: string;
  updated_at: string;
}

export interface CreateGaleryEventRequest {
  name: string;
  location: string;
  date: string; // ISO 8601 date string
  images_base64: string[];
}

export interface ModifyGaleryEventRequest {
  id: string;
  name: string;
  location: string;
  date: string; // ISO 8601 date string
  image_urls: string[];
  image_ids: string[];
}

/**
 * Get galery event by ID
 */
export async function getGaleryEventById(id: string): Promise<GaleryEvent> {
  return apiFetch<GaleryEvent>(`/galery_events/${id}`);
}

/**
 * List all galery events
 */

export async function listGaleryEvents(): Promise<GaleryEvent[]> {
  const res = await apiFetch<GaleryEvent[] | null>('/galery_events');
  return res ?? [];
}

/**
 * Create new galery event with images
 */
export async function createGaleryEvent(request: CreateGaleryEventRequest): Promise<GaleryEvent> {
  return apiFetch<GaleryEvent>('/galery_events', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}


export async function modifyGaleryEventRequest(request: ModifyGaleryEventRequest): Promise<GaleryEvent> {
  return apiFetch<GaleryEvent>('/galery_events', {
    method: 'PUT',
    body: JSON.stringify(request),
  });
}

// ==================
// TIMELINE API
// ==================

export interface TimelineEntry {
  id: string;
  name: string;
  text: string;
  location: string;
  date: string; // ISO 8601 date string
  created_at: string;
  updated_at: string;
}

export interface CreateTimelineEntryRequest {
  name: string;
  text: string;
  location: string;
  date: string; // ISO 8601 date string
}

export interface UpdateTimelineEntryRequest {
  name?: string;
  text?: string;
  location?: string;
  date?: string; // ISO 8601 date string
}

/**
 * Get timeline entry by ID
 */
export async function getTimelineEntryById(id: string): Promise<TimelineEntry> {
  return apiFetch<TimelineEntry>(`/timelineentries/${id}`);
}

/**
 * List all timeline entries
 */
export async function listTimelineEntries(): Promise<TimelineEntry[]> {
  return apiFetch<TimelineEntry[]>('/timelineentries');
}

/**
 * Create new timeline entry
 */
export async function createTimelineEntry(request: CreateTimelineEntryRequest): Promise<TimelineEntry> {
  return apiFetch<TimelineEntry>('/timelineentries', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}

/**
 * Update timeline entry
 */
export async function updateTimelineEntry(id: string, request: UpdateTimelineEntryRequest): Promise<TimelineEntry> {
  return apiFetch<TimelineEntry>(`/timelineentries/${id}`, {
    method: 'PUT',
    body: JSON.stringify(request),
  });
}

/**
 * Delete timeline entry
 */
export async function deleteTimelineEntry(id: string): Promise<void> {
  await apiFetchNoJSON(`/timelineentries/${id}`, {
    method: 'DELETE',
  });
}

// ==================
// TEXT API
// ==================

export interface Text {
  id: string;
  slug: string;
  content: string;
  page_slug?: string;
  page_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTextRequest {
  slug: string;
  content: string;
  page_slug?: string;
  page_id?: string;
}

export interface UpdateTextRequest {
  content?: string;
  page_slug?: string;
  page_id?: string;
}

/**
 * Get text by slug
 */
export async function getTextBySlug(slug: string): Promise<Text> {
  return apiFetch<Text>(`/texts/${slug}`);
}

/**
 * Get text by ID
 */
export async function getTextById(id: string): Promise<Text> {
  return apiFetch<Text>(`/texts/id/${id}`);
}

/**
 * Get texts by page slug
 */
export async function getTextsByPageSlug(pageSlug: string): Promise<Text[]> {
  return apiFetch<Text[]>(`/texts/page/slug/${pageSlug}`);
}

/**
 * List all texts
 */
export async function listTexts(): Promise<Text[]> {
  return apiFetch<Text[]>('/texts');
}

/**
 * Create new text
 */
export async function createText(request: CreateTextRequest): Promise<Text> {
  return apiFetch<Text>('/texts', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}

/**
 * Update text
 */
export async function updateText(id: string, request: UpdateTextRequest): Promise<Text> {
  return apiFetch<Text>(`/texts/${id}`, {
    method: 'PUT',
    body: JSON.stringify(request),
  });
}

/**
 * Delete text
 */
export async function deleteText(id: string): Promise<void> {
  await apiFetchNoJSON(`/texts/${id}`, {
    method: 'DELETE',
  });
}

// ==================
// EVENTS API (External)
// ==================

export interface Event {
  id: string;
  name: string;
  description?: string;
  starts_at: string;
  ends_at: string;
  location_name?: string;
  logo_url?: string;
  thumbnail_image_url?: string;
  link?: string;
}

/**
 * Get events with optional query parameters
 */
export async function getEvents(params?: {
  limit?: number;
  orderBy?: string;
  desc?: boolean;
}): Promise<Event[]> {
  const queryParams = new URLSearchParams();
  if (params?.limit) queryParams.append('limit', params.limit.toString());
  if (params?.orderBy) queryParams.append('orderBy', params.orderBy);
  if (params?.desc !== undefined) queryParams.append('desc', params.desc.toString());

  const query = queryParams.toString();
  const endpoint = `/events${query ? `?${query}` : ''}`;

  return apiFetch<Event[]>(endpoint);
}

// ==================
// UTILITY FUNCTIONS
// ==================

/**
 * Convert File to Base64 string
 */
export async function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      const result = reader.result as string;
      // Remove data URL prefix (e.g., "data:image/png;base64,")
      const base64 = result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = error => reject(error);
  });
}

/**
 * Convert multiple Files to Base64 strings
 */
export async function filesToBase64(files: File[]): Promise<string[]> {
  return Promise.all(files.map(file => fileToBase64(file)));
}
