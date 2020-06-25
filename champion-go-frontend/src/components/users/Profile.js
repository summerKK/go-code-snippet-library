import {useDispatch, useSelector} from "react-redux"
import React, {useState} from "react"
import _ from "lodash"
import Default from "../../assets/default.jpeg"
import {Redirect} from "react-router"
import Navigation from "../Navigation"
import {
    CardBody,
    Col,
    CustomInput,
    FormGroup,
    Row,
    Button,
    Form,
    Input,
    InputGroupAddon,
    InputGroup,
    Label
} from "reactstrap"
import Message from "../utils/Message"
import {updateUserAvatar, updateUser} from "../../store/modules/auth/actions/authActions"
import "./Profile.css"

const Profile = () => {

    const authState = useSelector(state => state.Auth)

    const dispatch = useDispatch()

    const [modal, setModal] = useState(false)
    const toggle = e => {
        setModal(!modal)
    }

    const authId = authState.currentUser ? authState.currentUser : ""

    const [uploadedFile, setUploadedFile] = useState()

    const [file, setFile] = useState()

    const [user, setUser] = useState({
        email: authState.currentUser.email,
        current_password: "",
        new_password: "",
    })

    const handleChange = e => {
        setUser({
            ...user,
            [e.target.name]: e.target.value,
        })
    }

    const handleImageChange = e => {
        let reader = new FileReader()
        let theFile = e.target.files[0]

        reader.onload = () => {
            setFile(theFile)
            setUploadedFile(reader.result)
            // 上传头像
            const formData = new FormData()
            formData.append("avatar", theFile)
            dispatch(updateUserAvatar(formData))
        }
        reader.readAsDataURL(theFile)
    }

    const submitUserAvatar = e => {
        e.preventDefault()
    }

    const clearUpInput = () => {
        setUser({
            ...user,
            current_password: "",
            new_password: "",
        })
    }

    const submitUser = e => {
        e.preventDefault()
        dispatch(updateUser({
            email: user.email,
            current_password: user.current_password,
            new_password: user.new_password,
        }, clearUpInput))
    }

    if (!authState.isAuthenticated) {
        return <Redirect to="/login"/>
    }

    let imagePreview = null
    if (_.get(authState, "currentUser.avatar_path") && !uploadedFile) {
        imagePreview = (<img className="img_style" src={_.get(authState, "currentUser.avatar_path")} alt="profile"/>)
    } else if (uploadedFile) {
        imagePreview = (<img className="img_style" src={uploadedFile} alt="profile"/>)
    } else {
        imagePreview = (<img className="img_style" src={Default} alt="profile"/>)
    }

    return (
        <React.Fragment>
            <div>
                <Navigation/>
            </div>
            <div className="post-style container">
                <div className="card-style">
                    <div className="text-center">
                        <h4>Update Profile</h4>
                    </div>
                    <Row className="mt-1">
                        <Col sm="12" md={{size: 10, offset: 1}}>
                            <FormGroup>
                                {
                                    authState.authSuccessImage && authState.avatarError === null
                                        ? (<Message msg={authState.authSuccessImage}/>)
                                        : ("")
                                }
                            </FormGroup>
                        </Col>
                    </Row>
                    <CardBody>
                        <div className="mb-3 text-center">
                            {imagePreview}
                        </div>
                        <Form onSubmit={submitUserAvatar} encType="multipart/form-data">
                            <Row col="12">
                                <InputGroup>
                                    <CustomInput
                                        id="user_avatar"
                                        type="file"
                                        accept="image/*"
                                        onChange={handleImageChange}
                                        dataBrowse={authState.isLoadingAvatar ? "Uploading" : "Upload"}
                                    />
                                    {authState.avatarError && authState.avatarError.Too_large ? (
                                        <small className="color-red">{authState.avatarError.Too_large}</small>
                                    ) : (
                                        ""
                                    )}
                                    {authState.avatarError && authState.avatarError.Not_Image ? (
                                        <small className="color-red">{authState.avatarError.Not_Image}</small>
                                    ) : (
                                        ""
                                    )}
                                </InputGroup>
                            </Row>
                        </Form>

                        <Form onSubmit={submitUser}>
                            <Row col="12">
                                <InputGroup className="mt-4">
                                    <InputGroupAddon addonType="prepend">UserName</InputGroupAddon>
                                    <Input
                                        type="text"
                                        name="username"
                                        value={authState.currentUser.username || ""}
                                        disabled
                                    />
                                </InputGroup>

                                <InputGroup className="mt-4">
                                    <InputGroupAddon addonType="prepend">Email</InputGroupAddon>
                                    <Input
                                        onChange={handleChange}
                                        type="text"
                                        name="email"
                                        value={user.email || ""}
                                    />
                                    {authState.userError && authState.userError.Required_email ? (
                                        <small className="color-red">{authState.userError.Required_email}</small>
                                    ) : (
                                        ""
                                    )}
                                    {authState.userError && authState.userError.Invalid_email ? (
                                        <small className="color-red">{authState.userError.Invalid_email}</small>
                                    ) : (
                                        ""
                                    )}
                                    {authState.userError && authState.userError.Taken_email ? (
                                        <small className="color-red">{authState.userError.Taken_email}</small>
                                    ) : (
                                        ""
                                    )}
                                </InputGroup>

                                <InputGroup className="mt-4">
                                    <InputGroupAddon addonType="prepend">Current Password</InputGroupAddon>
                                    <Input
                                        onChange={handleChange}
                                        type="password"
                                        name="current_password"
                                        value={user.current_password || ""}
                                    />
                                    {authState.userError && authState.userError.Password_mismatch ? (
                                        <small
                                            className="color-red">{authState.userError.Password_mismatch}</small>
                                    ) : (
                                        ""
                                    )}
                                    {authState.userError && authState.userError.Empty_current ? (
                                        <small className="color-red">{authState.userError.Empty_current}</small>
                                    ) : (
                                        ""
                                    )}
                                </InputGroup>

                                <InputGroup className="mt-4">
                                    <InputGroupAddon addonType="prepend">New Password</InputGroupAddon>
                                    <Input
                                        onChange={handleChange}
                                        type="password"
                                        name="new_password"
                                        value={user.new_password || ""}
                                    />
                                    {authState.userError && authState.userError.Invalid_password ? (
                                        <small
                                            className="color-red">{authState.userError.Invalid_password}</small>
                                    ) : (
                                        ""
                                    )}
                                    {authState.userError && authState.userError.Empty_new ? (
                                        <small className="color-red">{authState.userError.Empty_new}</small>
                                    ) : (
                                        ""
                                    )}
                                </InputGroup>
                            </Row>

                            <Row className="mt-4">
                                <Col sm="12" md={{size: 10, offset: 1}}>
                                    <FormGroup>
                                        {authState.authSuccessUser != null && authState.userError == null ? (
                                            <Message msg={authState.authSuccessUser}/>
                                        ) : (
                                            ""
                                        )}
                                    </FormGroup>
                                </Col>
                            </Row>

                            <Row className="mt-3" sm="12">
                                <Col xs="4" md={{offset: 4}}>
                                    {authState.isUpdatingUser ? (
                                        <Button
                                            color="primary"
                                            type="submit"
                                            block
                                            disabled
                                        >
                                            Updating...
                                        </Button>
                                    ) : (
                                        <Button
                                            color="primary"
                                            type="submit"
                                            block
                                        >
                                            Update
                                        </Button>
                                    )}
                                </Col>
                            </Row>
                        </Form>
                    </CardBody>
                </div>
            </div>
        </React.Fragment>
    )
}

export default Profile