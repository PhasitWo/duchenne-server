### DOCTOR WEBSITE
POST /login ok

GET /authState ok

GET /profile ok
PUT /profile ok

GET /doctor ok
POST /doctor ok createDoctorPermission
GET /doctor/:id ok
PUT /doctor/:id ok updateDoctorPermission
DELETE /doctor/:id ok deleteDoctorPermission

GET /patient ok
POST /patient ok createPatientPermission
GET /patient/:id ok
PUT /patient/:id ok updatePatientPermission
DELETE /patient/:id ok deletePatientPermission

GET /appointment? ok

GET /question? ok
GET /question/:id ok
PUT /question/:id/answer ok

### DOCTOR ROLES
user = []
admin = [createPatientPermission, updatePatientPermission, deletePatientPermission]
root = [{all of admin's permissions}, createDoctorPermission, updateDoctorPermission, deleteDoctorPermission]