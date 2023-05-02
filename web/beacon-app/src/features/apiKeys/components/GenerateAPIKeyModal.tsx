/* eslint-disable prettier/prettier */
import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikProvider, useFormik } from 'formik';
import { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { useParams } from 'react-router-dom';
import styled from 'styled-components';

import { useCreateProjectAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';
import { APIKeyDTO, NewAPIKey } from '@/features/apiKeys/types/createApiKeyService';
import { useFetchPermissions } from '@/hooks/useFetchPermissions';
import { useOrgStore } from '@/store';

import generateAPIKeyValidationSchema from '../schemas/generateAPIKeyValidationSchema';

type GenerateAPIKeyModalProps = {
  open: boolean;
  onSetKey: React.Dispatch<React.SetStateAction<any>>;
  onClose: () => void;
  onSetModalOpen?: React.Dispatch<React.SetStateAction<boolean>>;
  projectId?: string;
};

function GenerateAPIKeyModal({ open, onSetKey, onClose, projectId }: GenerateAPIKeyModalProps) {
  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const [fullSelected, setFullSelected] = useState(true);
  const [customSelected, setCustomSelected] = useState(false);
  const org = useOrgStore.getState() as any;
  const { permissions } = useFetchPermissions();
  const { createProjectNewKey, key, wasKeyCreated, isCreatingKey, hasKeyFailed, error } =
    useCreateProjectAPIKey();
  const handleCreateKey = ({ name, permissions }: any) => {
    const projID = projectId || (projectID as string);

    const payload = {
      projectID: projID,
      name,
      permissions,
    } satisfies APIKeyDTO;

    createProjectNewKey(payload);
  };
  if (wasKeyCreated) {
    onSetKey(key);
  }

  const formik = useFormik<NewAPIKey>({
    initialValues: {
      name: '',
      permissions: [''],
    },
    validationSchema: generateAPIKeyValidationSchema,
    onSubmit: (values) => {
      values.permissions = values.permissions.filter(Boolean);

      handleCreateKey(values as APIKeyDTO);
    },
  });

  const { values, setFieldValue, resetForm } = formik;

  useEffect(() => {
    if (fullSelected) {
      setFieldValue('permissions', permissions);
      setCustomSelected(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fullSelected, permissions]);

  useEffect(() => {
    if (customSelected) {
      setFieldValue('permissions', []);
      setFullSelected(false);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [customSelected]);

  useEffect(() => {
    if (wasKeyCreated) {
      onClose();
      resetForm();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [wasKeyCreated]);

  useEffect(() => {
    if (hasKeyFailed) {
      toast.error(`${(error as any)?.response?.data?.error}`, {
        id: 'create-api-key-error',
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hasKeyFailed]);

  return (
    <Modal
      open={open}
      title={<h1>Generate API Key for {org?.project?.name} project.</h1>}
      containerClassName="max-h-[90vh] overflow-scroll max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
      onClose={onClose}
      data-testid="keyModal"
    >
      <>
        {/* <button onClick={onClose} className="bg-transparent absolute top-4 right-4 border-none">
          <CloseIcon className="h-4 w-4" />
        </button> */}
        <FormikProvider value={formik}>
          <div>
            <p className="mb-5">Name your key and select access permissions.</p>
            <Form className="space-y-6">
              <fieldset>
                <h2 className="mb-3 font-semibold">Key Name</h2>
                <TextField
                  placeholder="default"
                  fullWidth
                  {...formik.getFieldProps('name')}
                  data-testid="keyName"
                />
                <ErrorMessage name="name" component="small" className="text-xs text-danger-500" />
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
                          setFieldValue('permissions', isSelected ? permissions : []);

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
                          setCustomSelected(!!isSelected);
                          // reset permissions
                          setFieldValue('permissions', []);
                        }}
                        isSelected={!!customSelected}
                      >
                        Check to grant access for each action.
                      </Checkbox>
                    </StyledFieldset>
                    {customSelected && (
                      <div className="mt-5 ml-5 w-full space-y-1 md:ml-10 md:w-1/2">
                        {permissions &&
                          permissions.length > 0 &&
                          permissions.map((permission: string, key: number) => (
                            <StyledFieldset key={key}>
                              <Checkbox
                                onChange={(isSelected) => {
                                  setFieldValue(
                                    'permissions',
                                    isSelected
                                      ? [...values.permissions, permission]
                                      : values.permissions.filter((p) => p !== permission)
                                  );
                                }}
                                isSelected={
                                  customSelected && values.permissions.includes(permission)
                                }
                              >
                                {permission}
                              </Checkbox>
                            </StyledFieldset>
                          ))}
                      </div>
                    )}
                  </Box>
                </div>
              </fieldset>
              <div className="item-center flex  justify-center">
                <Button
                  isLoading={isCreatingKey}
                  className="bg-[#6DD19C] px-6 py-3 font-semibold"
                  data-testid="generateKey"
                >
                  Generate API Key
                </Button>
              </div>
            </Form>
          </div>
        </FormikProvider>
      </>
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
