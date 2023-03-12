import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
import { Form, FormikProvider, useFormik } from 'formik';
import { useState } from 'react';
import styled from 'styled-components';

import { APIKeyDTO, NewAPIKey } from '@/features/apiKeys/types/createApiKeyService';
type GenerateAPIKeyModalProps = {
  data: any;
  open: boolean;
  isLoading: boolean;
  onCloseModal: () => void;
  onSuccessfulCreate: boolean;
  onCreateNewKey: ({ name, permission }: any) => void;
  setOpenGenerateAPIKeyModal: React.Dispatch<React.SetStateAction<boolean>>;
};
// if selected, then all permissions are selected
// if not selected, then all permissions are not selected

function GenerateAPIKeyModal({
  open,
  //setOpenGenerateAPIKeyModal,
  onSuccessfulCreate,
  data,
  onCloseModal,
  isLoading,
  onCreateNewKey,
}: GenerateAPIKeyModalProps) {
  const formik = useFormik<NewAPIKey>({
    initialValues: {
      name: '',
      permissions: [''],
    },
    onSubmit: (values) => {
      console.log('values', values);
      onCreateNewKey(values as APIKeyDTO);
    },
    // validationSchema: NewAPIKEYSchema,
  });
  const { values, setFieldValue, resetForm } = formik;
  const [fullSelected, setFullSelected] = useState(true);
  const [customSelected, setCustomSelected] = useState(false);

  if (onSuccessfulCreate) {
    resetForm();
  }

  return (
    <Modal
      open={open}
      title={<h1>Generate Your API Key</h1>}
      size="medium"
      onClose={onCloseModal}
      containerClassName="overflow-scroll h-[80vh] lg:h-[90vh]"
    >
      <FormikProvider value={formik}>
        <div>
          <p className="mb-5">
            Name your key and customize permissions. Or stick with the default.
          </p>
          <Form className="space-y-6">
            <fieldset>
              <h2 className="mb-3 font-semibold">Key Name</h2>
              <TextField placeholder="default" fullWidth {...formik.getFieldProps('name')} />
            </fieldset>
            <fieldset>
              <h2 className="mb-3 font-semibold">Permissions</h2>
              <div className="space-y-8">
                <Box>
                  <h2 className="mb-1 font-semibold">Full Access</h2>
                  <StyledFieldset>
                    <Checkbox
                      {...formik.getFieldProps('full')}
                      onChange={(isSelected) => {
                        setFullSelected(!!isSelected);
                        setCustomSelected(false);
                        // reset permissions
                        setFieldValue('permissions', []);
                        setFieldValue('permissions', isSelected ? data : []);

                        // clear all permissions and set full access
                      }}
                      isSelected={fullSelected}
                    >
                      Full Access (default) - Publish to topic, Subscribe to topic, Create Topic,
                      Read Topic, Delete Topic, Destroy Topic.
                    </Checkbox>
                  </StyledFieldset>
                </Box>
                <Box>
                  <h2 className="mb-1 font-semibold">Custom Access</h2>
                  <StyledFieldset>
                    <Checkbox
                      {...formik.getFieldProps('custom')}
                      onChange={(isSelected) => {
                        setFullSelected(false);
                        setCustomSelected(!!isSelected);
                        // reset permissions
                        setFieldValue('permissions', []);
                      }}
                      isSelected={!!customSelected}
                    >
                      Check to grant access for each action.
                    </Checkbox>
                  </StyledFieldset>
                  <div className="mt-5 ml-5 w-full space-y-1 md:ml-10 md:w-1/2">
                    {data &&
                      data.length > 0 &&
                      data.map((permission: string, key: number) => (
                        <StyledFieldset key={key}>
                          <Checkbox
                            onChange={(isSelected) => {
                              setFieldValue(
                                'permissions',
                                fullSelected ? [] : [...values.permissions]
                              );
                              setFullSelected(false);
                              setCustomSelected(isSelected);
                              setFieldValue(
                                'permissions',
                                isSelected
                                  ? [...values.permissions, permission]
                                  : values.permissions.filter((p) => p !== permission)
                              );
                            }}
                            isSelected={customSelected && values.permissions.includes(permission)}
                          >
                            {permission}
                          </Checkbox>
                        </StyledFieldset>
                      ))}
                  </div>
                </Box>
              </div>
            </fieldset>
            <div className="flex items-center justify-between">
              <p className="pl-5">
                Generate API Key for <span className="font-semibold">[insert project name]</span>{' '}
                project.
              </p>
              <Button isLoading={isLoading} className="bg-[#6DD19C] px-6 py-3 font-semibold">
                Generate API Key
              </Button>
            </div>
          </Form>
        </div>
      </FormikProvider>
    </Modal>
  );
}

export default GenerateAPIKeyModal;

const Box = styled.div`
  padding: 16px 20px;
  border: 1px solid gray;
  border-radius: 4px;
`;

const StyledFieldset = styled.fieldset`
  label {
    flex-direction: row-reverse !important;
    justify-content: space-between;
    font-size: 14px;
  }
  label svg {
    min-width: 25px;
  }
`;
