const API_URL = 'http://localhost:8000/api';


export const login = async (studentId: string, password: string) => {
    console.log("Attempting login with:", { studentId, password });
    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ studentId, password }),
            credentials: 'same-origin',
            mode: 'cors'
        });

        if (!response.ok) {
            let errorData;
            try {
                errorData = await response.json();
            } catch (e) {
                throw new Error(`Failed to login: ${response.status}`);
            }
            throw new Error(errorData.error || 'Failed to login');
        }

        return response.json();
    } catch (error) {
        console.error("Failed to login:", error);
        throw error;
    }
};

export const register = async (studentId: string, password: string, name: string, email?: string) => {
    const response = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ studentId, password, name, email }),
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to register: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to register');
    }

    return response.json();
};

export const fetchAvailableCourts = async (date: string, startTime: string, endTime: string) => {
    const response = await fetch(
        `${API_URL}/courts/available?date=${date}&startTime=${startTime}&endTime=${endTime}`
    );

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to fetch available courts: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to fetch available courts');
    }

    return response.json();
};

export const fetchUserBookings = async (token: string) => {
    const response = await fetch(`${API_URL}/bookings`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to fetch bookings: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to fetch bookings');
    }

    return response.json();
};

export const createBooking = async (token: string, bookingData: any) => {
    const response = await fetch(`${API_URL}/bookings`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(bookingData),
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to create booking: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to create booking');
    }

    return response.json();
};

export const cancelBooking = async (token: string, bookingId: string) => {
    const response = await fetch(`${API_URL}/bookings/${bookingId}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to cancel booking: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to cancel booking');
    }

    return response.json();
};

export const getProfile = async (token: string) => {
    const response = await fetch(`${API_URL}/profile`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to fetch profile: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to fetch profile');
    }

    return response.json();
};

export const updateProfile = async (token: string, name: string, email?: string) => {
    const response = await fetch(`${API_URL}/profile`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ name, email }),
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to update profile: ${response.status}`);
        }
        throw new Error(errorData.error || 'Failed to update profile');
    }

    return response.json();
};

export const uploadProfilePicture = async (token: string, file: File) => {
    const formData = new FormData();
    formData.append("profilePicture", file);

    const response = await fetch(`${API_URL}/profile/upload`, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
        },
        body: formData,
    });

    if (!response.ok) {
        let errorData;
        try {
            errorData = await response.json();
        } catch (e) {
            throw new Error(`Failed to upload profile picture: ${response.status}`);
        }
        throw new Error(errorData.error || "Failed to upload profile picture");
    }

    return response.json();
};

