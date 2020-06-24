import API_ROUTE from "../../../../apiRoute"
import axios from "axios"
import {history} from "../../../../history"
import setAuthorizationToken from "../../../../authorization/authorization"
import {
    SIGNUP_SUCCESS,
    SIGNUP_ERROR,
    LOGIN_SUCCESS,
    LOGOUT_SUCCESS,
    LOGIN_ERROR,
    SET_CURRENT_USER,
    CREATE_POST_SUCCESS,
    CREATE_POST_ERROR,
    CREATE_POST,
    FETCH_POSTS,
    UPDATE_USER_SUCCESS,
    UPDATE_USER_ERROR,
    UPDATE_USER_AVATAR,
    UPDATE_USER_AVATAR_ERROR,
    BEFORE_STATE,
    BEFORE_AVATAR_STATE,
    BEFORE_USER_STATE,
    FORGOT_PASSWORD_SUCCESS,
    FORGOT_PASSWORD_ERROR,
    RESET_PASSWORD_SUCCESS,
    RESET_PASSWORD_ERROR,
    DELETE_USER_SUCCESS,
    DELETE_USER_ERROR
} from "../authTypes"

export const signIn = (credentials) => {

    return async dispatch => {
        dispatch({type: BEFORE_STATE})
        try {
            let response = await axios.post(`${API_ROUTE}/login`, credentials)
            let userData = response.data.response
            window.localStorage.setItem("token", userData.token)
            window.localStorage.setItem("user_data", JSON.stringify(userData))
            setAuthorizationToken(userData.token)
            dispatch({type: LOGIN_SUCCESS, payload: userData})
        } catch (e) {
            dispatch({type: LOGIN_ERROR, payload: e.response.data.err})
        }
    }
}

export const logout = () => {

    return dispatch => {
        window.localStorage.removeItem("token")
        setAuthorizationToken(false)
        dispatch({type: LOGOUT_SUCCESS})
        window.localStorage.clear()
        history.push("/login")
    }
}

export const signUp = (newUser) => {
    return async dispatch => {
        dispatch({type: BEFORE_STATE})
        try {
            let response = await axios.post(`${API_ROUTE}/users`, newUser)
            dispatch({type: SIGNUP_SUCCESS})
            history.push("/login")
        } catch (e) {
            dispatch({type: SIGNUP_ERROR, payload: e.response.data.err})
        }
    }
}

export const updateUserAvatar = (updateUserAvatar) => {
    return async (dispatch, getState) => {
        dispatch({type: BEFORE_AVATAR_STATE})
        const {id} = getState().Auth.currentUser
        try {
            let response = await axios.put(`${API_ROUTE}/users/${id}`, updateUserAvatar, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })

            let updateUser = response.data.response
            window.localStorage.setItem("user_data", JSON.stringify(updateUser))

            dispatch({type: UPDATE_USER_AVATAR, payload: updateUser})
        } catch (e) {
            dispatch({type: UPDATE_USER_AVATAR_ERROR, payload: e.response.data.err})
        }
    }
}

export const updateUser = (updateUser, clearInput) => {

    return async (dispatch, getState) => {
        dispatch({type: BEFORE_USER_STATE})
        const {currentUser} = getState().Auth
        try {
            const res = await axios.put(`${API_ROUTE}/users/${currentUser.id}`, updateUser);
            let updatedUser = res.data.response

            dispatch({type: UPDATE_USER_SUCCESS, payload: updatedUser})
            window.localStorage.setItem('user_data', JSON.stringify(updatedUser)); //update the localstorages
            clearInput()
        } catch (err) {
            dispatch({type: UPDATE_USER_ERROR, payload: err.response.data.error})
        }
    }
}

export const deleteUser = (id) => {

    return async dispatch => {
        dispatch({type: BEFORE_STATE})
        try {
            const res = await axios.delete(`${API_ROUTE}/users/${id}`);
            let deleteMessage = res.data.response
            dispatch({type: DELETE_USER_SUCCESS, payload: deleteMessage})
            window.localStorage.clear(); //update the localstorage
            window.location.href = "/"
        } catch (err) {
            dispatch({type: DELETE_USER_ERROR, payload: err.response.data.error})
        }
    }
}

export const forgotPassword = (userEmail, clearInput) => {

    return async (dispatch) => {

        dispatch({type: BEFORE_STATE})

        try {
            const res = await axios.post(`${API_ROUTE}/password/forgot`, userEmail);
            let passwordRequest = res.data.response
            dispatch({type: FORGOT_PASSWORD_SUCCESS, payload: passwordRequest})
            clearInput()
        } catch (err) {
            dispatch({type: FORGOT_PASSWORD_ERROR, payload: err.response.data.error})
        }
    }
}

export const resetPassword = (details, clearInput) => {

    return async (dispatch) => {

        dispatch({type: BEFORE_STATE})

        try {
            const res = await axios.post(`${API_ROUTE}/password/reset`, details);
            let passwordRequest = res.data.response
            dispatch({type: RESET_PASSWORD_SUCCESS, payload: passwordRequest})
            clearInput()
        } catch (err) {
            dispatch({type: RESET_PASSWORD_ERROR, payload: err.response.data.error})
        }
    }
}