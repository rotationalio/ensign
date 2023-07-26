// /* eslint-disable prettier/prettier */
// import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
// import { ErrorMessage, Form, FormikProvider, useFormik } from 'formik';
// import generateAPIKeyValidationSchema from '../schemas/generateAPIKeyValidationSchema';
// import { NewAPIKey } from '@/features/apiKeys/types/createApiKeyService';

// interface GenerateAPIKeyModalProps {
//     handleSubmit: (values: NewAPIKey) => void;
//     initialValues: NewAPIKey;
// }
// const GenerateAPIKeyForm: React.FC<GenerateAPIKeyModalProps> = ({
//   handleSubmit,
//   initialValues,
// }) => {

//   const formik = useFormik<NewAPIKey>({
//     initialValues,
//     validationSchema: generateAPIKeyValidationSchema,
//       onSubmit: (values) => {
//           handleSubmit(values);
//         },
//   });
//   const { values, setFieldValue } = formik;
//   return (
//     <FormikProvider value={formik}>
//       <div>
//         <p className="mb-5">Name your key and select access permissions.</p>
//         <Form className="space-y-6">
//           <fieldset>
//             <h2 className="mb-3 font-semibold">Key Name</h2>
//             <TextField
//               placeholder="default"
//               fullWidth
//               {...formik.getFieldProps('name')}
//               data-testid="keyName"
//             />
//             <ErrorMessage name="name" component="small" className="text-xs text-danger-500" />
//           </fieldset>
//           <fieldset>
//             <h2 className="mb-3 font-semibold">Permissions</h2>
//             <div className="space-y-8">
//               <Box>
//                 <h2 className="mb-1 font-semibold">Full Access</h2>
//                 <StyledFieldset>
//                   <Checkbox
//                     {...formik.getFieldProps('full')}
//                     onChange={FullSelectHanlder}
//                     isSelected={fullSelected}
//                   >
//                     Full Access (default) - Publish to topic, Subscribe to topic, Create Topic, Read
//                     Topic, Delete Topic, Destroy Topic.
//                   </Checkbox>
//                 </StyledFieldset>
//               </Box>
//               <Box>
//                 <h2 className="mb-1 font-semibold">Custom Access</h2>
//                 <StyledFieldset>
//                   <Checkbox
//                     {...formik.getFieldProps('custom')}
//                     onChange={(isSelected) => {
//                       setCustomSelected(!!isSelected);
//                       // reset permissions
//                       setFieldValue('permissions', []);
//                     }}
//                     isSelected={!!customSelected}
//                   >
//                     Check to grant access for each action.
//                   </Checkbox>
//                 </StyledFieldset>
//                 {customSelected && (
//                   <div className="mt-5 ml-5 w-full space-y-1 md:ml-10 md:w-1/2">
//                     {permissions &&
//                       permissions.length > 0 &&
//                       permissions.map((permission: string, key: number) => (
//                         <StyledFieldset key={key}>
//                           <Checkbox
//                             onChange={(isSelected) => {
//                               setFieldValue(
//                                 'permissions',
//                                 isSelected
//                                   ? [...values.permissions, permission]
//                                   : values.permissions.filter((p) => p !== permission)
//                               );
//                             }}
//                             isSelected={customSelected && values.permissions.includes(permission)}
//                           >
//                             {permission}
//                           </Checkbox>
//                         </StyledFieldset>
//                       ))}
//                   </div>
//                 )}
//               </Box>
//             </div>
//           </fieldset>
//           <div className="item-center flex  justify-center">
//             <Button isLoading={isCreatingKey} data-testid="generateKey">
//               Generate API Key
//             </Button>
//           </div>
//         </Form>
//       </div>
//     </FormikProvider>
//   );
// };

// export default GenerateAPIKeyForm;

export {};
