/* eslint-disable prettier/prettier */
import { t,Trans } from '@lingui/macro';
import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikProvider, useFormik } from 'formik';
import { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { useParams } from 'react-router-dom';
import styled from 'styled-components';

import { useCreateProjectAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';
import { APIKeyDTO, NewAPIKey } from '@/features/apiKeys/types/createApiKeyService';
import { useFetchPermissions } from '@/hooks/useFetchPermissions';

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

  const formik = useFormik<NewAPIKey>({
    initialValues: {
      name: '',
      // set permissions to full access by default
      permissions: [],
    },
    validationSchema: generateAPIKeyValidationSchema,
    onSubmit: (values) => {
      values.permissions = values.permissions.filter(Boolean);

      handleCreateKey(values as APIKeyDTO);
    },
  });

  const FullSelectHandler = (isSelected: boolean) => {
    setFullSelected(!!isSelected);
    setCustomSelected(false);
    // reset permissions
    //setFieldValue('permissions', []);
    setFieldValue('permissions', isSelected ? permissions : []);

    // clear all permissions and set full access
  };

  const { values, setFieldValue, resetForm } = formik;

  useEffect(() => {
    if (fullSelected) {
      setFieldValue('permissions', permissions);
      setCustomSelected(false);
    }
  }, [fullSelected, permissions, setFieldValue]);

  useEffect(() => {
    if (customSelected) {
      setFieldValue('permissions', []);
      setFullSelected(false);
    }
  }, [customSelected, setFieldValue]);

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

  useEffect(() => {
    if (wasKeyCreated) {
      onSetKey(key);
    }
  }, [wasKeyCreated, key, onSetKey]);

  return (
    <Modal
      open={open}
      title={t`Customize Your API Key`}
      containerClassName="create-key-modal max-h-screen"
      onClose={onClose}
      data-testid="keyModal"
    >
      <>
        <FormikProvider value={formik}>
          <div>
            <p className="mb-3">
              <Trans>Name your key and select access permissions. We recommend that you:</Trans></p>
            <ul className="mb-5 ml-5 list-disc list-outside">
              <li><Trans>Use a prefix that identifies the application, service, or use case that the API key is for</Trans></li>
              <li><Trans>Use a unique identifier for each API key (e.g. random number, a timestamp, or a UUID)</Trans></li>
              <li><Trans>Use a descriptive suffix that indicates the permissions</Trans></li>
            </ul>
            <p className="mb-3"><Trans>Examples:</Trans></p>
            <ul className="mb-5 ml-5 list-disc list-outside">
              <li>my-app-12345-full-access</li>
              <li>llama2-20241120-read-only</li>
              <li>finance-team-uuid-write-only</li>
            </ul>
            <Form className="space-y-6">
              <fieldset>
                <h2 className="mb-3 font-semibold"><Trans>Key Name (required)</Trans></h2>
                <TextField
                  placeholder={t`Enter key name`}
                  fullWidth
                  {...formik.getFieldProps('name')}
                  data-testid="keyName"
                />
                <ErrorMessage name="name" component="small" className="text-xs text-danger-700" />
              </fieldset>
              <fieldset>
                <h2 className="mb-3 font-semibold"><Trans>Permissions</Trans></h2>
                <div className="space-y-8">
                  <Box>
                    <h2 className="mb-1 font-semibold"><Trans>Full Access</Trans></h2>
                    <StyledFieldset>
                      <Checkbox
                        {...formik.getFieldProps('full')}
                        onChange={FullSelectHandler}
                        isSelected={fullSelected}
                      >
                        <Trans>Full Access (default) - Publish to topic, Subscribe to topic, Create Topic,
                        Read Topic, Delete Topic, Destroy Topic.</Trans>
                      </Checkbox>
                    </StyledFieldset>
                  </Box>
                  <Box>
                    <h2 className="mb-1 font-semibold"><Trans>Custom Access</Trans></h2>
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
                        <Trans>Check to grant access for each action.</Trans>
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
                                <Trans>{permission}</Trans>
                              </Checkbox>
                            </StyledFieldset>
                          ))}
                      </div>
                    )}
                  </Box>
                </div>
              </fieldset>
              <div className="item-center flex  justify-center">
                <Button isLoading={isCreatingKey} data-testid="generateKey">
                  <Trans>Generate API Key</Trans>
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
